package spentcalories

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
)

// Основные константы, необходимые для расчетов.
const (
	lenStep                    = 0.65 // средняя длина шага.
	mInKm                      = 1000 // количество метров в километре.
	minInH                     = 60   // количество минут в часе.
	stepLengthCoefficient      = 0.45 // коэффициент для расчета длины шага на основе роста.
	walkingCaloriesCoefficient = 0.5  // коэффициент для расчета калорий при ходьбе
)

// parseTraining принимает строку с данными формата "3456,Ходьба,3h00m",
// которая содержит количество шагов, вид активности и продолжительность активности
// и возвращает четыре значения:
// int — количество шагов.
// string — вид активности.
// time.Duration — продолжительность активности.
// error — ошибку, если что-то пошло не так.
func parseTraining(data string) (int, string, time.Duration, error) {
	parts := strings.Split(data, ",")
	if len(parts) != 3 {
		return 0, "", 0, fmt.Errorf("incorrect data \"%s\" in parseTraining()", data)
	}
	num, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, "", 0, fmt.Errorf("steps parsing error: %s in parseTraining()", err)
	}
	if num <= 0 {
		return 0, "", 0, fmt.Errorf("steps number equal or below zero. data: %s in parseTraining", data)
	}
	activity := parts[1]
	dur, err := time.ParseDuration(parts[2])
	if err != nil {
		return 0, "", 0, fmt.Errorf("duration parsing error: %s in parseTraining()", err)
	}
	if dur <= 0 {
		return 0, "", 0, fmt.Errorf("duration equal or below zero. data: %s in parseTraining()", data)
	}
	return num, activity, dur, nil
}

// distance принимает количество шагов и рост пользователя в метрах, а возвращает дистанцию в километрах.
func distance(steps int, height float64) float64 {
	stepLength := height * stepLengthCoefficient
	distanceInM := float64(steps) * stepLength
	return distanceInM / mInKm
}

// meanSpeed принимает количество шагов steps, рост пользователя height
// и продолжительность активности duration  и возвращает среднюю скорость.
func meanSpeed(steps int, height float64, duration time.Duration) float64 {
	if duration <= 0 {
		return 0
	}
	distance := distance(steps, height)
	return distance / duration.Hours()
}

// TrainingInfo принимает:
// data string — строку с данными формата "3456,Ходьба,3h00m", которая содержит количество шагов, вид активности и продолжительность активности.
// weight, height float64 — вес (кг.) и рост (м.) пользователя.
// И возвращает два значения:
// string — строка с информацией о тренировке в формате, приведенном ниже.
// error — ошибку, при ее возникновении внутри функции.
// Пример возвращаемой строки:
// Тип тренировки: Бег
// Длительность: 0.75 ч.
// Дистанция: 10.00 км.
// Скорость: 13.34 км/ч
// Сожгли калорий: 18621.75
func TrainingInfo(data string, weight, height float64) (string, error) {
	steps, activityType, duration, err := parseTraining(data)
	if err != nil {
		log.Println(err)
		return "", err
	}
	var (
		durationInH float64 = duration.Hours()
		dist        float64 = distance(steps, height)
		speed       float64 = meanSpeed(steps, height, duration)
		calories    float64
	)
	switch activityType {
	case "Ходьба":
		calories, err = WalkingSpentCalories(steps, weight, height, duration)
		if err != nil {
			return "", err
		}
	case "Бег":
		calories, err = RunningSpentCalories(steps, weight, height, duration)
		if err != nil {
			return "", err
		}
	default:
		return "", errors.New("неизвестный тип тренировки")
	}
	return fmt.Sprintf("Тип тренировки: %s\nДлительность: %.2f ч.\nДистанция: %.2f км.\nСкорость: %.2f км/ч\nСожгли калорий: %.2f\n",
		activityType, durationInH, dist, speed, calories), nil
}

// RunningSpentCalories принимает:
// steps int — количество шагов.
// weight, height float64 — вес(кг.) и рост(м.) пользователя.
// duration time.Duration — продолжительность бега.
// И возвращает два значения:
// float64 — количество калорий, потраченных при беге.
// error — ошибку, если входные параметры некорректны (подумайте, какие значения параметров имеют смысл).
func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if duration <= 0 || steps <= 0 || weight <= 0 || height <= 0 {
		return 0, fmt.Errorf("incorrect data \"%v, %v, %v, %v\" in RunningSpentCalories()",
			steps, weight, height, duration)
	}
	meanSpeed := meanSpeed(steps, height, duration)
	durationInMinutes := duration.Minutes()
	return (weight * meanSpeed * durationInMinutes) / minInH, nil
}

// WalkingSpentCalories по сигнатуре совпадает с сигнатурой RunningSpentCalories().
// Она принимает:
// steps int — количество шагов.
// weight, height float64 — вес(кг.) и рост(м.) пользователя.
// duration time.Duration — продолжительность ходьбы.
// И возвращает два значения:
// float64 — количество калорий, потраченных при ходьбе.
// error — ошибку, если входные параметры некорректны (подумайте, какие значения параметров имеют смысл).
func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if duration <= 0 || steps <= 0 || weight <= 0 || height <= 0 {
		return 0, fmt.Errorf("incorrect data \"%v, %v, %v, %v\" in WalkingSpentCalories()",
			steps, weight, height, duration)
	}
	meanSpeed := meanSpeed(steps, height, duration)
	durationInMinutes := duration.Minutes()
	return (walkingCaloriesCoefficient * weight * meanSpeed * durationInMinutes) / minInH, nil
}
