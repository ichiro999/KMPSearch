package main

import (
	"bufio"
	"fmt"
	"os"
)

var next_map map[int]int

func Usage() {
	fmt.Printf("Use like:\n %s file match_string\n", os.Args[0])
}

func init_nextval(pattern_string string) (map[int]int, error) {
	if len(pattern_string) == 0 {
		return nil, fmt.Errorf("Empty patter string")
	}
	res := make(map[int]int)
	//utf-8
	r_str := []rune(pattern_string)

	res[1] = 0
	//char num: begin at 1
	var i int = int(1)
	var j int = int(0)
	for i < len(r_str) {
		if j == 0 || r_str[i-1] == r_str[j-1] {
			i += 1
			j += 1
			if r_str[i-1] != r_str[j-1] {
				res[i] = j
			} else {
				res[i] = res[j]
			}
		} else {
			j = res[j]
		}
	}
	return res, nil
}

func KMP_index(master_str, match_str string, pos uint32) (int32, error) {

	if len(master_str) == 0 || len(match_str) == 0 {
		return -1, fmt.Errorf("master string or match string is empty")
	}

	if int(pos) > len(master_str) {
		return -1, fmt.Errorf("the pos is illegal, it bigger than the len of master string")
	}
	if len(next_map) == 0 {
		return -1, fmt.Errorf("have not the next map")
	}

	if int(pos)+len(match_str) > len(master_str) {
		//no such match string at the master string
		return -1, nil
	}
	i := int(pos)
	j := int(1)

	for i <= len(master_str) && j <= len(match_str) {
		if j == 0 || master_str[i-1] == match_str[j-1] {
			i += 1
			j += 1
		} else {
			//next_map is the location of the char, like "abc", char 'a' is at 1
			j = next_map[j]
		}
	}
	if j >= len(match_str) {
		return int32(i - len(match_str)), nil
	} else {
		return -1, nil
	}
}

func main() {
	if len(os.Args) < 3 {
		Usage()
		return
	}
	file_name := os.Args[1]
	match_str := os.Args[2]

	var err error
	file, err := os.Open(file_name)
	if err != nil {
		fmt.Fprintf(os.Stdout, "error: %v", err)
		return
	}
	defer file.Close()

	next_map, err = init_nextval(match_str)
	if err != nil {
		fmt.Printf("init the next map error %v\n", err)
		return
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		s_str := scanner.Text()
		if len(s_str) == 0 {
			continue
		}
		idex, err := KMP_index(s_str, match_str, 1)
		if err != nil {
			fmt.Printf("KMP grep error %v\n", err)
			return
		}
		if idex == -1 {
			continue
		}

		fmt.Printf("%s have %s at %d\n", s_str, match_str, idex)
	}
	return
}
