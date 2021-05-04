/*
 * Copyright 2018 The openwallet Authors
 * This file is part of the openwallet library.
 *
 * The openwallet library is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Lesser General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * The openwallet library is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
 * GNU Lesser General Public License for more details.
 */

package tron

import (
	"testing"
)

func TestListWitnesses(t *testing.T) {
	if r, err := tw.ListWitnesses(); err != nil {
		t.Errorf("ListWitnesses failed: %v\n", err)
	} else {
		t.Logf("ListWitnesses return: \n\t%+v\n", r)
	}
}

func TestListNodes(t *testing.T) {
	if r, err := tw.ListNodes(); err != nil {
		t.Errorf("ListNodes failed: %v\n", err)
	} else {
		t.Logf("ListNodes return: \n\t%+v\n", r)
	}
}
