package manager

import (
	"fmt"
)

func AddNewScript(scriptName string) error {
	if err := validateScriptName(scriptName); err != nil {
		return err
	}

	fmt.Println(scriptName)

	return nil
}

func validateScriptName(name string) error {
	// Check length
	if len(name) == 0 {
		return fmt.Errorf("script name required")
	}
	if len(name) > 63 {
		return fmt.Errorf("script name must be 63 characters or less (got %d)", len(name))
	}

	// Check first character (must be lowercase letter, +, -)
	firstChar := rune(name[0])
	if !(firstChar == '+' || firstChar == '-' || (firstChar >= 'a' && firstChar <= 'z')) {
		return fmt.Errorf("script name must start with a lowercase letter or '+' (got '%c')", firstChar)
	}

	// Check all characters
	for i, ch := range name {
		isLowercase := ch >= 'a' && ch <= 'z'
		isDigit := ch >= '0' && ch <= '9'
		isHyphen := ch == '-'
		isUnderscore := ch == '_'
		isPlus := ch == '+'

		if !isLowercase && !isDigit && !isHyphen && !isUnderscore && !isPlus {
			return fmt.Errorf("script name contains invalid character '%c' (position %d): only lowercase letters, digits, hyphens, underscores, and plus signs are allowed", ch, i+1)
		}
	}

	return nil
}
