package lottus

import "fmt"

type MessageProcessor struct {
	Displayer Processor
	Process   Processor
}

type MessageProvider interface {
	Add(n string, m Message) error
	Get(n string) (error, Message)
}

type Processor func(r *Request, res Message) Message

func DefaultProcessor(app App) Processor {
	return func(req *Request, res Message) Message {
		err, val := res.Input.Validate(*req, res)

		if err != nil {
			res.Errors = append(res.Errors, fmt.Sprintf("%s", err))

			return res
		}

		_, msg := Redirect(app, val.NextMessage, req)

		return msg
	}
}
