// Lute - 一款对中文语境优化的 Markdown 引擎，支持 Go 和 JavaScript
// Copyright (c) 2019-present, b3log.org
//
// Lute is licensed under Mulan PSL v2.
// You can use this software according to the terms and conditions of the Mulan PSL v2.
// You may obtain a copy of Mulan PSL v2 at:
//         http://license.coscl.org.cn/MulanPSL2
// THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND, EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
// See the Mulan PSL v2 for more details.

package parse

import (
	"bytes"

	"github.com/riverzhou/lute/lex"
	"github.com/riverzhou/lute/util"
)

func (t *Tree) parseATXHeading() (ok bool, markers, content []byte, level int) {
	tokens := t.Context.currentLine[t.Context.nextNonspace:]
	var startCaret bool
	if (t.Context.ParseOption.VditorWYSIWYG || t.Context.ParseOption.VditorIR || t.Context.ParseOption.VditorSV) && bytes.HasPrefix(tokens, util.CaretTokens) {
		tokens = bytes.ReplaceAll(tokens, util.CaretTokens, nil)
		startCaret = true
	}

	marker := tokens[0]
	if lex.ItemCrosshatch != marker {
		return
	}

	var inCaret bool
	if (t.Context.ParseOption.VditorWYSIWYG || t.Context.ParseOption.VditorIR || t.Context.ParseOption.VditorSV) && bytes.Contains(tokens, []byte("#"+util.Caret+"#")) {
		tokens = bytes.ReplaceAll(tokens, util.CaretTokens, nil)
		inCaret = true
	}

	level = lex.Accept(tokens, lex.ItemCrosshatch)
	if 6 < level {
		return
	}

	var endCaret bool
	if (t.Context.ParseOption.VditorWYSIWYG || t.Context.ParseOption.VditorIR || t.Context.ParseOption.VditorSV) && bytes.HasPrefix(tokens[level:], util.CaretTokens) {
		tokens = bytes.ReplaceAll(tokens, util.CaretTokens, nil)
		endCaret = true
	}

	if level < len(tokens) && !lex.IsWhitespace(tokens[level]) {
		return
	}

	markers = t.Context.currentLine[t.Context.nextNonspace : t.Context.nextNonspace+level+1]

	content = make([]byte, 0, 256)
	_, tokens = lex.TrimLeft(tokens)
	_, tokens = lex.TrimLeft(tokens[level:])
	for _, token := range tokens {
		if lex.ItemNewline == token {
			break
		}
		content = append(content, token)
	}
	_, content = lex.TrimRight(content)
	closingCrosshatchIndex := len(content) - 1
	for ; 0 <= closingCrosshatchIndex; closingCrosshatchIndex-- {
		if lex.ItemCrosshatch == content[closingCrosshatchIndex] {
			continue
		}
		if lex.ItemSpace == content[closingCrosshatchIndex] {
			break
		} else {
			closingCrosshatchIndex = len(content)
			break
		}
	}

	if 0 >= closingCrosshatchIndex {
		content = make([]byte, 0, 0)
	} else if 0 < closingCrosshatchIndex {
		content = content[:closingCrosshatchIndex]
		_, content = lex.TrimRight(content)
	}

	if t.Context.ParseOption.VditorWYSIWYG || t.Context.ParseOption.VditorIR || t.Context.ParseOption.VditorSV {
		if startCaret || inCaret || endCaret {
			content = append(util.CaretTokens, content...)
		}

		if util.Caret == string(content) || "" == string(content) {
			return
		}
	}
	_, content = lex.TrimRight(content)
	ok = true
	return
}

func (t *Tree) parseSetextHeading() (level int) {
	ln := lex.TrimWhitespace(t.Context.currentLine)
	var caretInLn bool
	if t.Context.ParseOption.VditorWYSIWYG || t.Context.ParseOption.VditorIR || t.Context.ParseOption.VditorSV {
		if bytes.Contains(ln, util.CaretTokens) {
			caretInLn = true
			ln = bytes.ReplaceAll(ln, util.CaretTokens, nil)
			if 1 > len(ln) {
				return
			}
		}
	}

	start := 0
	marker := ln[start]
	if lex.ItemEqual != marker && lex.ItemHyphen != marker {
		return
	}

	length := len(ln)
	for ; start < length; start++ {
		token := ln[start]
		if lex.ItemEqual != token && lex.ItemHyphen != token {
			return
		}

		if 0 != marker {
			if marker != token {
				return
			}
		} else {
			marker = token
		}
	}

	level = 1
	if lex.ItemHyphen == marker {
		level = 2
	}

	if (t.Context.ParseOption.VditorWYSIWYG || t.Context.ParseOption.VditorIR || t.Context.ParseOption.VditorSV) && caretInLn {
		t.Context.oldtip.Tokens = lex.TrimWhitespace(t.Context.oldtip.Tokens)
		t.Context.oldtip.AppendTokens(util.CaretTokens)
	}
	return
}
