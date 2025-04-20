package lottus

import (
	"errors"
	"fmt"
)

type SelectionInput struct {
	Options       []Option
	ParameterName string
	NextMessage   *Message
}

func (s SelectionInput) Validate(r Request) (error, InputValidationResult) {
	var option Option
	for _, opt := range s.Options {
		if opt.Key == r.Prompt {
			option = opt
		}
	}

	if option.Key != r.Prompt {
		return errors.New(fmt.Sprintf("Could not find an option with key %s", r.Prompt)), InputValidationResult{}
	}

	msg := option.NextMessage

	if msg == nil {
		msg = s.NextMessage
	}

	if msg == nil {
		return errors.New(fmt.Sprintf("Could not find next message for message")), InputValidationResult{}
	}

	return nil, InputValidationResult{
		UserInput:   r.Prompt,
		Parameters:  []Parameter{{Name: s.ParameterName, Value: option.ParameterValue}},
		NextMessage: "",
	}
}

type TextInput struct {
	ParameterName string
	NextMessage   string
}

func (i TextInput) Validate(r Request, m Message) (error, InputValidationResult) {
	return nil, InputValidationResult{
		UserInput:   r.Prompt,
		NextMessage: i.NextMessage,
		Parameters:  []Parameter{{Name: i.ParameterName, Value: r.Prompt}},
	}
}

type NumberInput struct {
	NextMessage *Message
	Parameter   Parameter
}
