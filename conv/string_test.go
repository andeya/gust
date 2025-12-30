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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := QuoteJSONString(tt.input, tt.escapeHTML)
			assert.Equal(t, tt.expected, string(result), tt.description)
		})
	}
}
