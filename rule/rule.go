package foo

import (
  "fmt"
  "strings"
  "golang.org/x/net/html"
)

type Rule struct {
  Element string
  Id string
  Classes []string
}

const (
  MATCH = iota
  MISMATCH = iota
  NEUTRAL = iota //rule neither matches or doesn't match element
)

func Parse(rule_str string) (Rule, error) {
  terms := strings.Split(rule_str, ";")
  var rule Rule

  //first term will be the element type
  rule.Element = terms[0]
  
  for _, term := range terms[1:] {
    ops := strings.Split(term, "=")
    if len(ops) != 2 {
      return nil, errors.New("invalid operand")
    }

    if ops[0] != "class" && ops[0] != "id" {
      return nil, errors.New("invalid attribute (must be 'class' or 'id')")
    }
    
    op_arguments := strings.Split(ops[1], " ")
  }
}

func (Rule *rule) Match(token html.Token) bool {
  if t.Data != rule.Element {
    return false
  }

  id_matched := ""
  class_matched := ""
  rule_match := false

  for _, attr := range token.Attr {
    if attr.Key == "id" {
      if rule.Id != "" {
        if rule.Id != attr.Val {
          id_matched = "yes"
        }
      }
    } else if attr.Key == "class" {
      class_matched = "yes"
    }
  }
  
  if rule.Id != "" {
    if !id_matched {
      return false
    }
  }

  if len(rule.Classes != 0) {
    
  }
  //if rule.Id != "" && id_matched
}
