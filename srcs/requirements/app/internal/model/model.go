package model
import(
	"github.com/google/uuid"
)

type Wallet struct{
	Id 		uuid.UUID	`json:"valletId"`
	Amount 	int			`json:"amount"`
}