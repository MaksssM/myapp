package recommendations

import "fmt"

func GenerateRecommendations(userID string) []string {
	fmt.Println("Генерация рекомендаций для пользователя:", userID)
	return []string{"Post1", "Post2", "Post3"}
}
