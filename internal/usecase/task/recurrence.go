package task

import (
	"context"
	"fmt"
	"strings"
	"time"

	taskdomain "example.com/taskservice/internal/domain/task"
)

func (s *Service) StartRecurrenceWorker(ctx context.Context, interval time.Duration) {
	ticker := time.NewTicker(interval)
	go func() {
		for {
			select {
			case <-ctx.Done():
				ticker.Stop()
				return
			case <-ticker.C:
				if err := s.processRecurrences(ctx); err != nil {
					fmt.Println("recurrence worker error:", err)
				}
			}
		}
	}()
}

// processRecurrences создаёт новые задачи на основе периодичности
func (s *Service) processRecurrences(ctx context.Context) error {
	tasks, err := s.repo.List(ctx)
	if err != nil {
		return err
	}

	now := s.now()

	for _, t := range tasks {
		if t.RecurrenceType == "" {
			continue
		}

		// Проверяем, нужно ли создать новую задачу
		if shouldCreateTask(t, now) {
			newTask := &taskdomain.Task{
				Title:          t.Title,
				Description:    t.Description,
				Status:         taskdomain.StatusNew,
				RecurrenceType: t.RecurrenceType,
				RecurrenceRule: t.RecurrenceRule,
				CreatedAt:      now,
				UpdatedAt:      now,
			}

			_, err := s.repo.Create(ctx, newTask)
			if err != nil {
				fmt.Println("failed to create recurring task:", err)
			}
		}
	}

	return nil
}

// shouldCreateTask определяет, создаем ли задачу сегодня
func shouldCreateTask(t taskdomain.Task, now time.Time) bool {
	switch t.RecurrenceType {
	case taskdomain.Daily:
		// можно добавлять каждое n-й день, пока n=1
		return true
	case taskdomain.Monthly:
		// проверяем, есть ли день месяца в правилe
		if t.RecurrenceRule == "" {
			return false
		}
		// правило — число месяца, например "5" или "15,20"
		for _, day := range parseDays(t.RecurrenceRule) {
			if int(now.Day()) == day {
				return true
			}
		}
	case taskdomain.Specific:
		// правило — конкретные даты, например "2026-04-06,2026-04-07"
		for _, date := range parseDates(t.RecurrenceRule) {
			if now.Format("2006-01-02") == date.Format("2006-01-02") {
				return true
			}
		}
	case taskdomain.EvenDays:
		return now.Day()%2 == 0
	case taskdomain.OddDays:
		return now.Day()%2 != 0
	}

	return false
}

// parseDays возвращает срез дней из строки "1,5,10"
func parseDays(rule string) []int {
	var days []int
	for _, part := range splitAndTrim(rule) {
		var day int
		fmt.Sscanf(part, "%d", &day)
		if day >= 1 && day <= 30 {
			days = append(days, day)
		}
	}
	return days
}

// parseDates возвращает срез time.Time из строки "2026-04-06,2026-04-07"
func parseDates(rule string) []time.Time {
	var dates []time.Time
	for _, part := range splitAndTrim(rule) {
		if date, err := time.Parse("2006-01-02", part); err == nil {
			dates = append(dates, date)
		}
	}
	return dates
}

func splitAndTrim(s string) []string {
	parts := []string{}
	for _, p := range strings.Split(s, ",") {
		p = strings.TrimSpace(p)
		if p != "" {
			parts = append(parts, p)
		}
	}
	return parts
}
