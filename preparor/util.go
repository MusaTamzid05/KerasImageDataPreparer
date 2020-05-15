package preparor


func common(str1 , str2 string) string {

	matches := [] string {}

	for index , _ := range str1 {
		for index2 , _ := range str2 {
			if str1[index] == str2[index2] {

				if string(str1[index]) == " " {
					continue
				}
				
				i := index
				j := index2
				match := ""
				matchFound := true

				for matchFound {
					
					if j >= len(str2) || i >= len(str1) {
						matchFound = false

						if match != "" {
							matches = append(matches , match)
						}
						continue
					}

					if str1[i] != str2[j] {
						matches = append(matches , match)
						matchFound = false
						continue
					}

	
					match += string(str1[i])
					i += 1
					j += 1

				}
			}
		}
	}

	biggestMatch := ""
	max := -1

	for _ , match := range matches {
		if len(match)  >  max {
			max = len(match)
			biggestMatch = match
		}
	}

	return biggestMatch

}

