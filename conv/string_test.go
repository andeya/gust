package conv

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIndent(t *testing.T) {
	tests := []struct {
		name     string
		text     string
		prefix   string
		expected string
	}{
		{
			name:     "simple text",
			text:     "line1\nline2",
			prefix:   "  ",
			expected: "  line1\n  line2",
		},
		{
			name:     "with trailing newline",
			text:     "line1\nline2\n",
			prefix:   "  ",
			expected: "  line1\n  line2\n",
		},
		{
			name:     "empty prefix",
			text:     "line1\nline2",
			prefix:   "",
			expected: "line1\nline2",
		},
		{
			name:     "empty text",
			text:     "",
			prefix:   "  ",
			expected: "",
		},
		{
			name:     "empty text and prefix",
			text:     "",
			prefix:   "",
			expected: "",
		},
		{
			name:     "single line",
			text:     "single line",
			prefix:   "  ",
			expected: "  single line",
		},
		{
			name:     "multiple newlines",
			text:     "line1\n\nline2\n\n\nline3",
			prefix:   "  ",
			expected: "  line1\n\n  line2\n\n\n  line3",
		},
		{
			name:     "only newlines",
			text:     "\n\n\n",
			prefix:   "  ",
			expected: "\n\n\n",
		},
		{
			name:     "unicode characters",
			text:     "中文\n测试",
			prefix:   "  ",
			expected: "  中文\n  测试",
		},
		{
			name:     "tab prefix",
			text:     "line1\nline2",
			prefix:   "\t",
			expected: "\tline1\n\tline2",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Indent(tt.text, tt.prefix)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestToSnakeCase(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"XxYy", "XxYy", "xx_yy"},
		{"_XxYy", "_XxYy", "_xx_yy"},
		{"TcpRpc", "TcpRpc", "tcp_rpc"},
		{"ID", "ID", "id"},
		{"UserID", "UserID", "user_id"},
		{"RPC", "RPC", "rpc"},
		{"TCP_RPC", "TCP_RPC", "tcp_rpc"},
		{"wakeRPC", "wakeRPC", "wake_rpc"},
		{"_TCP__RPC", "_TCP__RPC", "_tcp__rpc"},
		{"_TcP__RpC_", "_TcP__RpC_", "_tc_p__rp_c_"},
		{"empty string", "", ""},
		{"single uppercase", "A", "a"},
		{"single lowercase", "a", "a"},
		{"all uppercase", "ABC", "abc"},
		{"all lowercase", "abc", "abc"},
		{"mixed with numbers", "UserID123", "user_id123"},
		{"already snake_case", "user_id", "user_id"},
		{"camelCase", "camelCase", "camel_case"},
		{"PascalCase", "PascalCase", "pascal_case"},
		{"XMLHttpRequest", "XMLHttpRequest", "xml_http_request"},
		{"HTTPSConnection", "HTTPSConnection", "https_connection"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ToSnakeCase(tt.input)
			assert.Equal(t, tt.expected, result, "input: %q", tt.input)
			// Should be idempotent
			result2 := ToSnakeCase(result)
			assert.Equal(t, tt.expected, result2, "should be idempotent for %q", result)
		})
	}
}

func TestToCamelCase(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"underscore only", "_", "_"},
		{"xx_yy", "xx_yy", "XxYy"},
		{"_xx_yy", "_xx_yy", "_XxYy"},
		{"id", "id", "Id"},
		{"user_id", "user_id", "UserId"},
		{"rpc", "rpc", "Rpc"},
		{"tcp_rpc", "tcp_rpc", "TcpRpc"},
		{"wake_rpc", "wake_rpc", "WakeRpc"},
		{"_tcp___rpc", "_tcp___rpc", "_Tcp__Rpc"},
		{"_tc_p__rp_c__", "_tc_p__rp_c__", "_TcP_RpC__"},
		{"empty string", "", ""},
		{"single char", "a", "A"},
		{"already camelCase", "CamelCase", "CamelCase"},
		{"with numbers", "user_id_123", "UserId123"},
		{"multiple underscores", "a__b__c", "A__B__C"},
		{"leading underscores", "___abc", "___Abc"},
		{"trailing underscores", "abc___", "Abc___"},
		// Additional test cases for 100% coverage
		{"underscore before uppercase", "user_Name", "User_Name"},
		{"underscore at end", "abc_", "Abc_"},
		{"underscore before non-letter", "abc_@def", "Abc_@def"},
		{"underscore before uppercase letter", "abc_Adef", "Abc_Adef"},
		{"single underscore after leading", "_a", "_A"},
		{"underscore before digit", "abc_123def", "Abc123def"},
		{"uppercase letter", "ABC", "ABC"},
		{"uppercase with underscore", "ABC_DEF", "ABC_DEF"},
		{"digit only", "123", "123"},
		{"underscore between digits", "123_456", "123456"},
		{"mixed case with underscore", "aBc_DeF", "ABc_DeF"},
		{"underscore before uppercase after leading", "_a_B", "_A_B"},
		{"complex pattern", "a_b_C__d", "AB_C_D"},
		{"underscore sequence with uppercase", "a__B", "A__B"},
		{"trailing underscore after letter", "abc_", "Abc_"},
		{"multiple underscores with uppercase", "a___B", "A__B"},
		{"underscore before uppercase in middle", "test_Value", "Test_Value"},
		{"single char uppercase", "A", "A"},
		{"uppercase after underscore", "_A", "_A"},
		{"letter after multiple underscores", "a___b", "A__B"},
		{"pattern with uppercase", "x_Y__Z", "X_Y_Z"},
		// Test uppercase letter after underscore (next >= 'A' && next <= 'Z')
		{"underscore before uppercase", "test_Value", "Test_Value"},
		{"underscore before uppercase after leading", "_a_B", "_A_B"},
		// Test last character is underscore (i == lastIdx)
		{"underscore at very end", "abc_", "Abc_"},
		// Test uppercase letter in hasLetterAfter check
		{"multiple underscores before uppercase", "a__B", "A__B"},
		{"underscore sequence with uppercase after", "test___Value", "Test__VAlue"},
		// Test uppercase prevChar in hasSingleUnderscoreBefore check
		{"uppercase letter before double underscore", "A__B", "A__B"},
		{"uppercase pattern", "X_Y__Z", "X_Y_Z"},
		// Test uppercase letter in main loop (b >= 'A' && b <= 'Z')
		{"uppercase without capitalize", "ABC", "ABC"},
		{"uppercase after underscore", "_ABC", "_ABC"},
		{"mixed case", "aBc", "ABc"},
		// Test edge cases
		{"single underscore at end", "a_", "A_"},
		{"underscore before non-letter char", "a_@b", "A_@b"},
		{"underscore before digit at end", "a_1", "A1"},
		{"multiple underscores at end", "a___", "A___"},
		{"underscore before uppercase at start", "_A", "_A"},
		{"complex uppercase pattern", "A_B__C", "A_B_C"},
		// Test i == leadingUnderscores case
		{"double underscore right after leading", "__a", "__A"},
		// Test non-letter prevChar case
		{"double underscore after digit", "1__a", "1__A"},
		// Test i-2 < leadingUnderscores case
		{"double underscore after single leading", "_a__b", "_A__B"},
		// Test i-3 < leadingUnderscores case
		{"double underscore after double leading", "__a__b", "__A__B"},
		// Test non-letter prevPrevPrev case
		{"double underscore after digit and underscore", "1_a__b", "1A__B"},
		// Test i == lastIdx case (last character is underscore)
		{"underscore at very end", "abc_", "Abc_"},
		// Test next >= 'A' && next <= 'Z' case (uppercase letter after underscore)
		{"underscore before uppercase", "test_Value", "Test_Value"},
		// Test underscoreCount == 1 in else branch (line 264)
		{"single underscore in else branch", "a_b_c", "ABC"},
		// Test underscoreCount == 1 case in else branch
		{"single underscore between words", "a_b", "AB"},
		// Test hasLetterAfter with uppercase
		{"double underscore before uppercase", "a__B", "A__B"},
		// Test underscoreCount < 2 case
		{"single underscore trailing", "a_", "A_"},
		// Test hasLetterAfter false with non-letter
		{"double underscore before non-letter", "a__@", "A__@"},
		// Test underscoreCount >= 3 with hasLetterAfter
		{"triple underscore before letter", "a___b", "A__B"},
		// Test underscoreCount == 2 with hasLetterAfter but no single underscore before
		{"double underscore after letter", "a__b", "A__B"},
		// Test i-2 == leadingUnderscores edge case
		{"double underscore after single leading underscore", "_a__b", "_A__B"},
		// Test next >= 'A' && next <= 'Z' case (line 156 else branch - uppercase after underscore)
		{"underscore before uppercase", "test_Value", "Test_Value"},
		// Test i == leadingUnderscores case (line 162 else branch)
		{"underscore right after leading with lowercase", "_a_b", "_AB"},
		// Test i-2 < leadingUnderscores case (line 224 else branch)
		{"double underscore with i-2 < leadingUnderscores", "_a__b", "_A__B"},
		// Test i-3 < leadingUnderscores case (line 245 else branch)
		{"double underscore with i-3 < leadingUnderscores", "__a__b", "__A__B"},
		// Test non-letter prevChar case (line 240 else branch)
		{"double underscore after non-letter", "1__a", "1__A"},
		// Test non-letter prevPrevPrev case (line 247 else branch)
		{"double underscore after digit underscore letter", "1_a__b", "1A__B"},
		// Test underscoreCount == 1 in else branch (line 264)
		{"single underscore in else branch", "a_b", "AB"},
		// Test i == leadingUnderscores in line 162 (underscore right after leading)
		{"underscore right after leading", "_a", "_A"},
		// Test i-2 == leadingUnderscores edge case (line 224)
		{"double underscore at boundary", "_a__b", "_A__B"},
		// Test i-3 == leadingUnderscores edge case (line 245)
		{"double underscore at boundary 2", "__a__b", "__A__B"},
		// Test prevChar is digit case (line 240)
		{"double underscore after digit letter", "1a__b", "1A__B"},
		// Test prevPrevPrev is digit case (line 247)
		{"double underscore after digit underscore letter pattern", "1_a__b", "1A__B"},
		// Test next is uppercase letter (not lowercase) after underscore
		{"underscore before uppercase letter", "test_Value", "Test_Value"},
		// Test underscoreCount == 1 in else branch (line 264) - single underscore
		{"single underscore in else branch", "a_b", "AB"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ToCamelCase(tt.input)
			assert.Equal(t, tt.expected, result, "input: %q", tt.input)
		})
	}
}

func TestToPascalCase(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"underscore only", "_", "_"},
		{"xx_yy", "xx_yy", "XxYy"},
		{"_xx_yy", "_xx_yy", "XxYy"},
		{"id", "id", "ID"},
		{"user_id", "user_id", "UserID"},
		{"rpc", "rpc", "RPC"},
		{"tcp_rpc", "tcp_rpc", "TCPRPC"},
		{"wake_rpc", "wake_rpc", "WakeRPC"},
		{"___tcp___rpc", "___tcp___rpc", "TCPRPC"},
		{"_tc_p__rp_c__", "_tc_p__rp_c__", "TcPRpC"},
		{"empty string", "", ""},
		{"single char", "a", "A"},
		{"api_key", "api_key", "APIKey"},
		{"http_request", "http_request", "HTTPRequest"},
		{"xml_parser", "xml_parser", "XMLParser"},
		{"json_data", "json_data", "JSONData"},
		{"url_path", "url_path", "URLPath"},
		{"with numbers", "user_id_123", "UserID123"},
		{"mixed case", "xml_http_request", "XMLHttpRequest"},
		// Additional test cases for 100% coverage
		{"initialism at start", "api_key", "APIKey"},
		{"initialism at end", "key_api", "KeyAPI"},
		{"initialism without next part", "xml_http", "XMLHTTP"},
		{"initialism with next initialism", "xml_http_json", "XMLHTTPJSON"},
		{"single char part", "a_b_c", "ABC"},
		{"empty parts in split", "a__b", "AB"},
		{"multiple empty parts", "a___b", "AB"},
		{"initialism followed by non-initialism", "tcp_request", "TCPRequest"},
		{"non-initialism followed by initialism", "request_tcp", "RequestTCP"},
		{"all initialisms", "tcp_rpc", "TCPRPC"},
		{"single char", "a", "A"},
		{"single initialism", "api", "API"},
		{"initialism at start with non-initialism", "api_request", "APIRequest"},
		// Test idx+1 >= len(parts) case (last part) - line 344 else branch
		{"initialism at end after initialism", "request_api", "RequestAPI"},
		{"initialism at end", "api_json", "APIJSON"}, // last part, idx+1 == len(parts)
		// Test commonInitialisms[nextPart] == true case - line 346 else branch
		{"initialism followed by initialism", "api_json", "APIJSON"},
		{"three initialisms", "api_json_xml", "APIJSONXML"}, // nextPart is also initialism
		{"initialism in middle with initialism after", "request_api_json", "RequestAPIJSON"},
		// Test single char part (len(part) == 1) - line 361 else branch
		{"single char parts", "a_b_c", "ABC"},
		{"single char with initialism", "a_id_b", "AIDB"},
		{"single char only", "a", "A"},
		{"single char with underscore", "a_b", "AB"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ToPascalCase(tt.input)
			assert.Equal(t, tt.expected, result, "input: %q", tt.input)
		})
	}
}

func TestDecodeHTMLEntities(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		radix    int
		expected string
	}{
		{
			name:     "hexadecimal entities",
			input:    `{"info":[["color","&#5496;&#5561;&#8272;&#7c;&#7eff;&#8272;"]]｝`,
			radix:    16,
			expected: `{"info":[["color","咖啡色|绿色"]]｝`,
		},
		{
			name:     "decimal entities",
			input:    `&#65;&#66;&#67;`,
			radix:    10,
			expected: `ABC`,
		},
		{
			name:     "no entities",
			input:    `hello world`,
			radix:    16,
			expected: `hello world`,
		},
		{
			name:     "empty string",
			input:    ``,
			radix:    16,
			expected: ``,
		},
		{
			name:     "invalid entity",
			input:    `&#invalid;`,
			radix:    16,
			expected: `&#invalid;`,
		},
		{
			name:     "mixed valid and invalid",
			input:    `&#65;&#invalid;&#67;`,
			radix:    10,
			expected: `A&#invalid;C`,
		},
		{
			name:     "entity at start",
			input:    `&#65;BC`,
			radix:    10,
			expected: `ABC`,
		},
		{
			name:     "entity at end",
			input:    `AB&#67;`,
			radix:    10,
			expected: `ABC`,
		},
		{
			name:     "multiple same entities",
			input:    `&#65;&#65;&#65;`,
			radix:    10,
			expected: `AAA`,
		},
		{
			name:     "base 8",
			input:    `&#101;&#102;&#103;`,
			radix:    8,
			expected: `ABC`,
		},
		{
			name:     "base 2",
			input:    `&#1000001;&#1000010;&#1000011;`,
			radix:    2,
			expected: `ABC`,
		},
		{
			name:     "unicode beyond ASCII",
			input:    `&#20013;&#25991;`,
			radix:    10,
			expected: `中文`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := DecodeHTMLEntities(tt.input, tt.radix)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestDecodeUnicodeEscapes(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		radix    int
		expected string
	}{
		{
			name:     "hexadecimal code points",
			input:    `{"info":[["color","\u5496\u5561\u8272\u7c\u7eff\u8272"]]｝`,
			radix:    16,
			expected: `{"info":[["color","咖啡色|绿色"]]｝`,
		},
		{
			name:     "no code points",
			input:    `hello world`,
			radix:    16,
			expected: `hello world`,
		},
		{
			name:     "empty string",
			input:    ``,
			radix:    16,
			expected: ``,
		},
		{
			name:     "code point at start",
			input:    `\u0041BC`,
			radix:    16,
			expected: `ABC`,
		},
		{
			name:     "code point at end",
			input:    `AB\u0043`,
			radix:    16,
			expected: `ABC`,
		},
		{
			name:     "invalid code point",
			input:    `\uinvalid`,
			radix:    16,
			expected: `\uinvalid`,
		},
		{
			name:     "long code point",
			input:    `\u12345678`,
			radix:    16,
			expected: "\u1234" + "5678", // First 4 chars parsed, rest kept
		},
		{
			name:     "multiple code points",
			input:    `\u0041\u0042\u0043`,
			radix:    16,
			expected: `ABC`,
		},
		{
			name:     "mixed valid and invalid",
			input:    `\u0041\uinvalid\u0043`,
			radix:    16,
			expected: `A\uinvalidC`,
		},
		{
			name:     "decimal radix",
			input:    `\u0065\u0066\u0067`,
			radix:    10,
			expected: `ABC`,
		},
		{
			name:     "empty escape sequence",
			input:    `\u`,
			radix:    16,
			expected: `\u`,
		},
		{
			name:     "short code point",
			input:    `\u41BC`,
			radix:    16,
			expected: `ABC`,
		},
		{
			name:     "unicode beyond BMP",
			input:    `\uD83D\uDE00`,
			radix:    16,
			expected: `😀`,
		},
		{
			name:     "escaped backslash before u",
			input:    `\\u0041`,
			radix:    16,
			expected: `\\u0041`, // Should not decode escaped backslash
		},
		{
			name:     "even number of backslashes",
			input:    `\\\\u0041`,
			radix:    16,
			expected: `\\\\u0041`, // Even number means escaped, not a real escape sequence
		},
		{
			name:     "escape at start",
			input:    `\u0041`,
			radix:    16,
			expected: `A`, // pos == 0 case
		},
		{
			name:     "long code point with invalid parse",
			input:    `\uGGGG1234`,
			radix:    16,
			expected: `\uGGGG1234`, // Invalid parse, keep original
		},
		{
			name:     "short code point",
			input:    `\u41`,
			radix:    16,
			expected: `A`, // Should pad to 4 chars
		},
		{
			name:     "code point with rest",
			input:    `\u0041BC`,
			radix:    16,
			expected: `ABC`, // Normal 4-char code point
		},
		{
			name:     "high surrogate without low",
			input:    `\uD83D`,
			radix:    16,
			expected: "\ufffd", // High surrogate without low surrogate, decoded as replacement char
		},
		{
			name:     "high surrogate with invalid low",
			input:    `\uD83D\u0041`,
			radix:    16,
			expected: "\ufffdA", // High surrogate with non-low surrogate, decoded as replacement char
		},
		{
			name:     "high surrogate with low but invalid parse",
			input:    `\uD83D\uGGGG`,
			radix:    16,
			expected: "\ufffd\\uGGGG", // Invalid parse of low surrogate, high decoded as replacement char
		},
		{
			name:     "high surrogate with low but out of bounds",
			input:    `\uD83D\uDE00`,
			radix:    16,
			expected: `😀`, // Valid surrogate pair
		},
		{
			name:     "code point with rest after parse",
			input:    `\u12345678`,
			radix:    16,
			expected: "\u1234" + "5678", // Long code point, parse first 4, keep rest
		},
		{
			name:     "invalid code point parse",
			input:    `\uGGGG`,
			radix:    16,
			expected: `\uGGGG`, // Invalid parse, keep original
		},
		{
			name:     "non-hex radix with special case",
			input:    `\u41BC`,
			radix:    10,
			expected: `\u41BC`, // Non-hex radix, don't apply special case
		},
		{
			name:     "code point with different first chars",
			input:    `\u51BC`,
			radix:    16,
			expected: "冼", // Different first chars, normal parse as 4-char code point
		},
		{
			name:     "empty code point at end",
			input:    `test\u`,
			radix:    16,
			expected: `test\u`, // Empty code point
		},
		{
			name:     "escape at start pos 0",
			input:    `\u0041`,
			radix:    16,
			expected: `A`, // pos == 0 case
		},
		{
			name:     "high surrogate with low but out of bounds",
			input:    `\uD83D\uDE00`,
			radix:    16,
			expected: `😀`, // Valid surrogate pair
		},
		{
			name:     "high surrogate with low but segmentEnd out of bounds",
			input:    `\uD83D\uDE`,
			radix:    16,
			expected: "\ufffd\u00DE", // segmentEnd+len(escapeSeq)+4 > len(str), high surrogate becomes replacement char
		},
		{
			name:     "high surrogate with low but invalid parse",
			input:    `\uD83D\uGGGG`,
			radix:    16,
			expected: "\ufffd\\uGGGG", // err2 != nil case
		},
		{
			name:     "high surrogate with low but not in range",
			input:    `\uD83D\u0041`,
			radix:    16,
			expected: "\ufffdA", // nextCodePoint not in 0xDC00-0xDFFF
		},
		{
			name:     "code point with rest",
			input:    `\u12345678`,
			radix:    16,
			expected: "\u1234" + "5678", // len(rest) > 0
		},
		{
			name:     "code point without rest",
			input:    `\u0041`,
			radix:    16,
			expected: `A`, // len(rest) == 0
		},
		{
			name:     "long code point with invalid parse",
			input:    `\uGGGG1234`,
			radix:    16,
			expected: `\uGGGG1234`, // err != nil for long code point
		},
		{
			name:     "high surrogate with nextIdx < 0",
			input:    `\uD83D`,
			radix:    16,
			expected: "\ufffd", // nextIdx < 0, so no low surrogate check
		},
		{
			name:     "high surrogate with segmentEnd out of bounds",
			input:    `\uD83D\uDE`,
			radix:    16,
			expected: "\ufffd\u00DE", // segmentEnd+len(escapeSeq)+4 > len(str)
		},
		{
			name:     "high surrogate with invalid low parse",
			input:    `\uD83D\uGGGG`,
			radix:    16,
			expected: "\ufffd\\uGGGG", // err2 != nil
		},
		{
			name:     "high surrogate with low not in range",
			input:    `\uD83D\u0041`,
			radix:    16,
			expected: "\ufffdA", // nextCodePoint not in 0xDC00-0xDFFF
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := DecodeUnicodeEscapes(tt.input, tt.radix)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestNormalizeWhitespace(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name: "multiple spaces",
			input: `# authenticate method 

		//  comment2	

		/*  some other 
			  comments */
		`,
			expected: `# authenticate method
	// comment2
	/* some other
	comments */
	`,
		},
		{
			name:     "empty string",
			input:    ``,
			expected: ``,
		},
		{
			name:     "single space",
			input:    `hello world`,
			expected: `hello world`,
		},
		{
			name:     "multiple newlines",
			input:    "line1\n\n\nline2",
			expected: "line1\n\nline2",
		},
		{
			name:     "multiple spaces between words",
			input:    "hello    world",
			expected: "hello world",
		},
		{
			name:     "tabs and spaces",
			input:    "hello\t\t\tworld",
			expected: "hello\tworld",
		},
		{
			name:     "mixed whitespace",
			input:    "hello \t \n \t world",
			expected: "hello world",
		},
		{
			name:     "trailing spaces",
			input:    "hello world   ",
			expected: "hello world",
		},
		{
			name:     "leading spaces",
			input:    "   hello world",
			expected: "hello world",
		},
		{
			name:     "only spaces",
			input:    "     ",
			expected: "",
		},
		{
			name:     "only newlines",
			input:    "\n\n\n",
			expected: "\n\n",
		},
		{
			name:     "space before newline",
			input:    "text \n\nmore",
			expected: "text\nmore",
		},
		{
			name:     "tab before newline",
			input:    "text\t\n\nmore",
			expected: "text\nmore",
		},
		{
			name:     "paragraph break preserved",
			input:    "para1\n\npara2",
			expected: "para1\n\npara2",
		},
		{
			name:     "windows line endings",
			input:    "line1\r\n\r\nline2",
			expected: "line1\n\nline2",
		},
		{
			name:     "mixed line endings",
			input:    "line1\r\n\nline2",
			expected: "line1\n\nline2",
		},
		{
			name:     "unicode whitespace",
			input:    "hello\u00A0\u00A0world", // Non-breaking spaces
			expected: "hello world",
		},
		{
			name:     "complex whitespace block",
			input:    "start   \t  \n  \t  \n  \t  end",
			expected: "start\nend",
		},
		// Additional test cases for 100% coverage
		{
			name:     "tab after newline",
			input:    "text\n\tindented",
			expected: "text\n\tindented",
		},
		{
			name:     "tab after newline with space before",
			input:    "text \n\tindented",
			expected: "text\n\tindented",
		},
		{
			name:     "multiple tabs after newline",
			input:    "text\n\t\tindented",
			expected: "text\n\tindented", // Multiple tabs normalized to single tab
		},
		{
			name:     "newline with unicode whitespace after",
			input:    "text\n\u00A0more",
			expected: "text\nmore", // Unicode whitespace after newline is normalized
		},
		{
			name:     "newline with unicode whitespace before",
			input:    "text\u00A0\nmore",
			expected: "text\nmore",
		},
		{
			name:     "invalid UTF-8 character",
			input:    "text\xFF\xFEmore",
			expected: "text\xFF\xFEmore",
		},
		{
			name:     "invalid UTF-8 with whitespace",
			input:    "text \xFE\xFF more",
			expected: "text \xFE\xFF more",
		},
		{
			name:     "tab not after newline",
			input:    "text\t\tmore",
			expected: "text\tmore",
		},
		{
			name:     "tab with newline count",
			input:    "text\n\n\tindented",
			expected: "text\n\n\tindented",
		},
		{
			name:     "space with newline count",
			input:    "text\n\n more",
			expected: "text\n\nmore", // Space after newline is normalized
		},
		{
			name:     "unicode whitespace with newline",
			input:    "text\n\u00A0\u00A0more",
			expected: "text\nmore", // Unicode whitespace after newline is normalized
		},
		{
			name:     "unicode whitespace without newline",
			input:    "text\u00A0\u00A0more",
			expected: "text more",
		},
		{
			name:     "carriage return only",
			input:    "text\rmore",
			expected: "text\nmore",
		},
		{
			name:     "carriage return with space before",
			input:    "text \rmore",
			expected: "text\nmore",
		},
		{
			name:     "carriage return with tab before",
			input:    "text\t\rmore",
			expected: "text\nmore",
		},
		{
			name:     "multiple carriage returns",
			input:    "text\r\r\rmore",
			expected: "text\n\nmore",
		},
		{
			name:     "carriage return newline",
			input:    "text\r\nmore",
			expected: "text\nmore",
		},
		{
			name:     "tab before newline with content",
			input:    "text\t\nmore",
			expected: "text\nmore",
		},
		{
			name:     "space before newline with content",
			input:    "text \nmore",
			expected: "text\nmore",
		},
		{
			name:     "newline with tab after (indentation)",
			input:    "text\n\t\tmore",
			expected: "text\n\tmore", // Multiple tabs normalized to single tab
		},
		{
			name:     "newline with space then tab",
			input:    "text\n \tmore",
			expected: "text\nmore",
		},
		{
			name:     "leading tab",
			input:    "\ttext",
			expected: "text",
		},
		{
			name:     "trailing tab after newline",
			input:    "text\n\t",
			expected: "text\n\t",
		},
		{
			name:     "trailing tab without newline",
			input:    "text\t",
			expected: "text",
		},
		{
			name:     "tab after newline in middle",
			input:    "start\n\tmiddle\n\tend",
			expected: "start\n\tmiddle\n\tend",
		},
		{
			name:     "unicode whitespace character",
			input:    "text\u2000\u2001more",
			expected: "text more",
		},
		{
			name:     "unicode whitespace with newline count",
			input:    "text\n\u2000more",
			expected: "text\nmore", // Unicode whitespace after newline is normalized
		},
		// Test firstNonSpaceAfter == '\t' case (line 697)
		{
			name:     "newline with tab after (indentation pattern)",
			input:    "text\n\tmore",
			expected: "text\n\tmore", // firstNonSpaceAfter == '\t', startsWithTab == true
		},
		// Test firstNonSpaceAfter != '\t' case (line 697)
		{
			name:     "newline with space then tab",
			input:    "text\n \tmore",
			expected: "text\nmore", // firstNonSpaceAfter != '\t', startsWithTab == false
		},
		// Test firstNonSpaceAfter == 0 case (line 697)
		{
			name:     "newline with only spaces after",
			input:    "text\n  more",
			expected: "text\nmore", // firstNonSpaceAfter == 0, startsWithTab == false
		},
		// Test firstCharAfter == '\t' case (line 706)
		{
			name:     "newline directly followed by tab",
			input:    "text\n\tmore",
			expected: "text\n\tmore", // firstCharAfter == '\t', startsWithTab == true
		},
		// Test firstCharAfter != '\t' case (line 706)
		{
			name:     "newline with space before tab",
			input:    "text\n \tmore",
			expected: "text\nmore", // firstCharAfter == ' ', startsWithTab == false
		},
		// Test i+1 >= len(input) case (line 703 else branch)
		{
			name:     "newline at end",
			input:    "text\n",
			expected: "text\n", // i+1 >= len(input), firstCharAfter == 0
		},
		// Test firstNonSpaceAfter from Unicode whitespace (line 676)
		{
			name:     "newline with unicode whitespace then char",
			input:    "text\n\u00A0more",
			expected: "text\nmore", // Unicode whitespace, firstNonSpaceAfter set to 'm'
		},
		// Test firstNonSpaceAfter from regular character (line 683)
		{
			name:     "newline with regular char",
			input:    "text\nmore",
			expected: "text\nmore", // Regular character, firstNonSpaceAfter set to 'm'
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := NormalizeWhitespace(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func ExampleQuoteJSONString() {
	str := `<>&{}""`
	fmt.Printf("%s\n", QuoteJSONString(str, true))
	fmt.Printf("%s\n", QuoteJSONString(str, false))
	// Output:
	// "\u003c\u003e\u0026{}\"\""
	// "<>&{}\"\""
}

func TestQuoteJSONString(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		escapeHTML  bool
		expected    string
		description string
	}{
		{
			name:        "escape HTML",
			input:       `<>&{}""`,
			escapeHTML:  true,
			expected:    `"\u003c\u003e\u0026{}\"\""`,
			description: "HTML characters should be escaped",
		},
		{
			name:        "no HTML escape",
			input:       `<>&{}""`,
			escapeHTML:  false,
			expected:    `"<>&{}\"\""`,
			description: "HTML characters should not be escaped",
		},
		{
			name:        "empty string",
			input:       ``,
			escapeHTML:  false,
			expected:    `""`,
			description: "empty string should return empty JSON string",
		},
		{
			name:        "newline and tab",
			input:       "hello\n\tworld",
			escapeHTML:  false,
			expected:    `"hello\n\tworld"`,
			description: "newline and tab should be escaped",
		},
		{
			name:        "backslash and quote",
			input:       `hello\"world`,
			escapeHTML:  false,
			expected:    `"hello\\\"world"`,
			description: "backslash and quote should be escaped",
		},
		{
			name:        "control characters",
			input:       "\x00\x01\x02",
			escapeHTML:  false,
			expected:    `"\u0000\u0001\u0002"`,
			description: "control characters should be escaped",
		},
		{
			name:        "unicode characters",
			input:       "中文",
			escapeHTML:  false,
			expected:    `"中文"`,
			description: "unicode characters should be preserved",
		},
		{
			name:        "line separator",
			input:       "\u2028",
			escapeHTML:  false,
			expected:    `"\u2028"`,
			description: "line separator should be escaped",
		},
		{
			name:        "paragraph separator",
			input:       "\u2029",
			escapeHTML:  false,
			expected:    `"\u2029"`,
			description: "paragraph separator should be escaped",
		},
		{
			name:        "invalid UTF-8",
			input:       "\xFF\xFE",
			escapeHTML:  false,
			expected:    `"\ufffd\ufffd"`,
			description: "invalid UTF-8 should be replaced with replacement character",
		},
		{
			name:        "mixed content",
			input:       "Hello <script>alert('XSS')</script>",
			escapeHTML:  true,
			expected:    `"Hello \u003cscript\u003ealert('XSS')\u003c/script\u003e"`,
			description: "mixed content with HTML should be escaped when escapeHTML=true",
		},
		{
			name:        "carriage return",
			input:       "hello\rworld",
			escapeHTML:  false,
			expected:    `"hello\rworld"`,
			description: "carriage return should be escaped",
		},
		{
			name:        "form feed",
			input:       "hello\fworld",
			escapeHTML:  false,
			expected:    `"hello\u000cworld"`,
			description: "form feed should be escaped",
		},
		{
			name:        "backspace",
			input:       "hello\bworld",
			escapeHTML:  false,
			expected:    `"hello\u0008world"`,
			description: "backspace should be escaped",
		},
		// Additional test cases for 100% coverage
		{
			name:        "writeStart equals i",
			input:       "\x00",
			escapeHTML:  false,
			expected:    `"\u0000"`,
			description: "writeStart equals i case",
		},
		{
			name:        "line separator",
			input:       "\u2028",
			escapeHTML:  false,
			expected:    `"\u2028"`,
			description: "line separator should be escaped",
		},
		{
			name:        "paragraph separator",
			input:       "\u2029",
			escapeHTML:  false,
			expected:    `"\u2029"`,
			description: "paragraph separator should be escaped",
		},
		{
			name:        "line separator with prefix",
			input:       "text\u2028more",
			escapeHTML:  false,
			expected:    `"text\u2028more"`,
			description: "line separator with prefix",
		},
		{
			name:        "paragraph separator with prefix",
			input:       "text\u2029more",
			escapeHTML:  false,
			expected:    `"text\u2029more"`,
			description: "paragraph separator with prefix",
		},
		{
			name:        "safe characters with escapeHTML true",
			input:       "safe<>",
			escapeHTML:  true,
			expected:    `"safe\u003c\u003e"`,
			description: "safe characters with escapeHTML true",
		},
		{
			name:        "safe characters with escapeHTML false",
			input:       "safe<>",
			escapeHTML:  false,
			expected:    `"safe<>"`,
			description: "safe characters with escapeHTML false",
		},
		{
			name:        "writeStart equals len at end",
			input:       "test",
			escapeHTML:  false,
			expected:    `"test"`,
			description: "writeStart equals len at end",
		},
		{
			name:        "multi-byte character",
			input:       "中文",
			escapeHTML:  false,
			expected:    `"中文"`,
			description: "multi-byte character",
		},
		{
			name:        "mixed ASCII and multi-byte",
			input:       "hello中文world",
			escapeHTML:  false,
			expected:    `"hello中文world"`,
			description: "mixed ASCII and multi-byte",
		},
		{
			name:        "invalid UTF-8 at start",
			input:       "\xFF\xFEtest",
			escapeHTML:  false,
			expected:    `"\ufffd\ufffdtest"`,
			description: "invalid UTF-8 at start",
		},
		{
			name:        "invalid UTF-8 in middle",
			input:       "test\xFF\xFEmore",
			escapeHTML:  false,
			expected:    `"test\ufffd\ufffdmore"`,
			description: "invalid UTF-8 in middle",
		},
		{
			name:        "invalid UTF-8 with prefix write",
			input:       "prefix\xFF\xFE",
			escapeHTML:  false,
			expected:    `"prefix\ufffd\ufffd"`,
			description: "invalid UTF-8 with prefix write",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := QuoteJSONString(tt.input, tt.escapeHTML)
			assert.Equal(t, tt.expected, string(result), tt.description)
		})
	}
}
