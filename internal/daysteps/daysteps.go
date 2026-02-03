package daysteps

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/Yandex-Practicum/tracker/internal/spentcalories"
)

const (
	// Длина одного шага в метрах
	stepLength = 0.65
	// Количество метров в одном километре
	mInKm = 1000
)

// parsePackage парсит строку формата "678,0h50m", где:
// 678 — количество шагов.
// 0h50m — продолжительность прогулки.
// Переводит эти данные в int и time.Duration соответственно и возвращает их.
func parsePackage(data string) (int, time.Duration, error) {
	parts := strings.Split(data, ",")
	if len(parts) != 2 {
		return 0, 0, fmt.Errorf("incorrect data \"%s\" in parsePackage()", data)
	}
	num, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, 0, fmt.Errorf("steps parsing error: %w in parsePackage()", err)
	}
	if num <= 0 {
		return 0, 0, fmt.Errorf("steps number equal or below zero. data: %s in parsePackage()", data)
	}
	dur, err := time.ParseDuration(parts[1])
	if err != nil {
		return 0, 0, fmt.Errorf("duraion parsing error: %w in parsePackage()", err)
	}
	if dur <= 0 {
		return 0, 0, fmt.Errorf("duration equal or below zero. data: %s in parsePackage()", data)
	}
	return num, dur, nil
}

// DayActionInfo парсит строку с данными с помощью parsePackage(),
// вычисляет дистанцию в километрах, количество потраченных калорий и возвращает строку в виде:
// Количество шагов: 792.
// Дистанция составила 0.51 км.
// Вы сожгли 221.33 ккал.
func DayActionInfo(data string, weight, height float64) string {
	numOfSteps, walkDuration, err := parsePackage(data)
	if err != nil {
		log.Println(err)
		return ""
	}
	if numOfSteps < 0 {
		return ""
	}
	distanceInKm := float64(numOfSteps) * stepLength / mInKm
	calories, err := spentcalories.WalkingSpentCalories(numOfSteps, weight, height, walkDuration)
	if err != nil {
		log.Println(err)
		return ""
	}
	return fmt.Sprintf("Количество шагов: %v.\nДистанция составила %.2f км.\nВы сожгли %.2f ккал.\n",
		numOfSteps, distanceInKm, calories)
}
