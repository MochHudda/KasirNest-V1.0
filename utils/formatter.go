package utils

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// FormatCurrency formats number as Rupiah currency
func FormatCurrency(amount float64) string {
	// Convert to string with 2 decimal places
	str := fmt.Sprintf("%.2f", amount)

	// Split integer and decimal parts
	parts := strings.Split(str, ".")
	intPart := parts[0]
	decPart := parts[1]

	// Add thousand separators to integer part
	intPartFormatted := AddThousandSeparators(intPart)

	// Return formatted currency
	if decPart == "00" {
		return "Rp " + intPartFormatted
	}
	return "Rp " + intPartFormatted + "," + decPart
}

// AddThousandSeparators adds thousand separators to a number string
func AddThousandSeparators(numStr string) string {
	// Remove any existing separators
	numStr = strings.ReplaceAll(numStr, ".", "")
	numStr = strings.ReplaceAll(numStr, ",", "")

	// Add dots as thousand separators (Indonesian format)
	n := len(numStr)
	if n <= 3 {
		return numStr
	}

	var result strings.Builder
	for i, digit := range numStr {
		if i > 0 && (n-i)%3 == 0 {
			result.WriteString(".")
		}
		result.WriteRune(digit)
	}

	return result.String()
}

// FormatNumber formats integer with thousand separators
func FormatNumber(num int) string {
	return AddThousandSeparators(strconv.Itoa(num))
}

// FormatDate formats time as Indonesian date format
func FormatDate(t time.Time) string {
	months := []string{
		"Januari", "Februari", "Maret", "April", "Mei", "Juni",
		"Juli", "Agustus", "September", "Oktober", "November", "Desember",
	}

	return fmt.Sprintf("%d %s %d", t.Day(), months[t.Month()-1], t.Year())
}

// FormatDateTime formats time as Indonesian date and time format
func FormatDateTime(t time.Time) string {
	return fmt.Sprintf("%s %02d:%02d", FormatDate(t), t.Hour(), t.Minute())
}

// FormatTime formats time as HH:MM
func FormatTime(t time.Time) string {
	return fmt.Sprintf("%02d:%02d", t.Hour(), t.Minute())
}

// FormatDateShort formats date in short format (DD/MM/YYYY)
func FormatDateShort(t time.Time) string {
	return fmt.Sprintf("%02d/%02d/%d", t.Day(), t.Month(), t.Year())
}

// FormatDateTimeShort formats date and time in short format
func FormatDateTimeShort(t time.Time) string {
	return fmt.Sprintf("%s %s", FormatDateShort(t), FormatTime(t))
}

// ParseCurrency parses currency string to float64
func ParseCurrency(currencyStr string) (float64, error) {
	// Remove currency symbol and spaces
	cleaned := strings.TrimSpace(currencyStr)
	cleaned = strings.TrimPrefix(cleaned, "Rp")
	cleaned = strings.TrimPrefix(cleaned, "rp")
	cleaned = strings.TrimSpace(cleaned)

	// Remove thousand separators (dots)
	cleaned = strings.ReplaceAll(cleaned, ".", "")

	// Replace decimal comma with dot
	cleaned = strings.ReplaceAll(cleaned, ",", ".")

	// Parse to float64
	return strconv.ParseFloat(cleaned, 64)
}

// FormatPercentage formats number as percentage
func FormatPercentage(value float64) string {
	return fmt.Sprintf("%.1f%%", value)
}

// TruncateString truncates string to specified length with ellipsis
func TruncateString(str string, maxLength int) string {
	if len(str) <= maxLength {
		return str
	}

	if maxLength <= 3 {
		return str[:maxLength]
	}

	return str[:maxLength-3] + "..."
}

// FormatFileSize formats file size in bytes to human readable format
func FormatFileSize(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}

	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}

	units := []string{"KB", "MB", "GB", "TB"}
	return fmt.Sprintf("%.1f %s", float64(bytes)/float64(div), units[exp])
}

// CleanString removes extra whitespaces and trims string
func CleanString(str string) string {
	// Replace multiple spaces with single space
	re := strings.NewReplacer("  ", " ", "\t", " ", "\n", " ", "\r", "")
	cleaned := re.Replace(str)

	// Trim leading and trailing spaces
	return strings.TrimSpace(cleaned)
}
