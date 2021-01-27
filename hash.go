package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

type Payment struct {
	PayerINN          string `json:"payerinn"`         // ИНН плательщика
	PayerKPP          string `json:"payerkpp"`         // КПП плательщика
	PayeeINN          string `json:"payeeinn"`         // ИНН получателя
	PayeeKPP          string `json:"payeekpp"`         // КПП получателя
	PayeeBIK          string `json:"payeebik"`         // БИК получателя
	PayeeCheckAccount string `json:"payeecheckccount"` // расч. счет получателя
	PayeeCorrAccount  string `json:"payeecorraccount"` // корр. счет получателя
	Amount            string `json:"amount"`           // сумма платежа
	Details           string `json:"details"`          // назначение платежа (не более 20 символов)
}

func generateHash(p Payment) string {
	details := string(([]rune(p.Details))[:20])
	hash := sha256.Sum256([]byte(p.PayerINN + `|` + p.PayerKPP + `|` + p.PayeeINN + `|` + p.PayeeKPP + `|` + p.PayeeBIK + `|` + p.PayeeCheckAccount + `|` + p.PayeeCorrAccount + `|` + p.Amount + `|` + details))
	return hex.EncodeToString(hash[:])
}

func main() {
	var p Payment
	p.PayerINN = "3528000xxx"
	p.PayerKPP = "997550xxx"
	p.PayeeINN = "3528000xxx"
	p.PayeeKPP = "997550xxx"
	p.PayeeBIK = "044525xxx"
	p.PayeeCheckAccount = "30101810400000000xxx"
	p.PayeeCorrAccount = "40702810500020106xxx"
	p.Amount = "1000"
	p.Details = "Пополнение счета на текущие расходы НДС не облагается!"

	fmt.Print(generateHash(p))
}
