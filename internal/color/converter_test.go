package color

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewConverter_ValidInputs(t *testing.T) {
	t.Run("should_parse_hex6_with_hash", func(t *testing.T) {
		converter, err := NewConverter("#ffffff")

		require.NoError(t, err)
		require.NotNil(t, converter)
		assert.Equal(t, "#ffffff", converter.value)
	})

	t.Run("should_parse_hex6_uppercase", func(t *testing.T) {
		converter, err := NewConverter("#FFFFFF")

		require.NoError(t, err)
		assert.Equal(t, "#ffffff", converter.value)
	})

	t.Run("should_parse_hex6_without_hash", func(t *testing.T) {
		converter, err := NewConverter("ffffff")

		require.NoError(t, err)
		assert.Equal(t, "#ffffff", converter.value)
	})

	t.Run("should_parse_hex3_with_hash", func(t *testing.T) {
		converter, err := NewConverter("#fff")

		require.NoError(t, err)
		assert.Equal(t, "#ffffff", converter.value)
	})

	t.Run("should_parse_hex3_without_hash", func(t *testing.T) {
		converter, err := NewConverter("fff")

		require.NoError(t, err)
		assert.Equal(t, "#ffffff", converter.value)
	})

	t.Run("should_parse_0x_format", func(t *testing.T) {
		converter, err := NewConverter("0xffffff")

		require.NoError(t, err)
		assert.Equal(t, "#ffffff", converter.value)
	})

	t.Run("should_parse_0x_uppercase", func(t *testing.T) {
		converter, err := NewConverter("0XFFFFFF")

		require.NoError(t, err)
		assert.Equal(t, "#ffffff", converter.value)
	})

	t.Run("should_handle_mixed_case", func(t *testing.T) {
		converter, err := NewConverter("#Ab12eF")

		require.NoError(t, err)
		assert.Equal(t, "#ab12ef", converter.value)
	})
}

func TestNewConverter_InvalidInputs(t *testing.T) {
	t.Run("should_return_error_for_empty_string", func(t *testing.T) {
		converter, err := NewConverter("")

		assert.Error(t, err)
		assert.Nil(t, converter)
		assert.Contains(t, err.Error(), "empty")
	})

	t.Run("should_return_error_for_invalid_characters", func(t *testing.T) {
		converter, err := NewConverter("#gggggg")

		assert.Error(t, err)
		assert.Nil(t, converter)
		assert.Contains(t, err.Error(), "invalid")
	})

	t.Run("should_return_error_for_too_short_input", func(t *testing.T) {
		converter, err := NewConverter("#ff")

		assert.Error(t, err)
		assert.Nil(t, converter)
	})

	t.Run("should_return_error_for_too_long_input", func(t *testing.T) {
		converter, err := NewConverter("#fffffff")

		assert.Error(t, err)
		assert.Nil(t, converter)
	})

	t.Run("should_return_error_for_invalid_0x_format", func(t *testing.T) {
		converter, err := NewConverter("0xgggggg")

		assert.Error(t, err)
		assert.Nil(t, converter)
	})

	t.Run("should_return_error_for_just_hash", func(t *testing.T) {
		converter, err := NewConverter("#")

		assert.Error(t, err)
		assert.Nil(t, converter)
	})

	t.Run("should_return_error_for_random_text", func(t *testing.T) {
		converter, err := NewConverter("not_a_color")

		assert.Error(t, err)
		assert.Nil(t, converter)
	})
}

func TestNewConverter_EdgeCases(t *testing.T) {
	t.Run("should_handle_black_color", func(t *testing.T) {
		converter, err := NewConverter("#000000")

		require.NoError(t, err)
		assert.Equal(t, "#000000", converter.value)
	})

	t.Run("should_handle_white_color", func(t *testing.T) {
		converter, err := NewConverter("#ffffff")

		require.NoError(t, err)
		assert.Equal(t, "#ffffff", converter.value)
	})

	t.Run("should_expand_hex3_correctly", func(t *testing.T) {
		converter, err := NewConverter("#abc")

		require.NoError(t, err)
		assert.Equal(t, "#aabbcc", converter.value)
	})
}

func TestConverter_ToHex6(t *testing.T) {
	t.Run("should_return_hex6_format", func(t *testing.T) {
		converter, err := NewConverter("#ab12ef")
		require.NoError(t, err)

		result := converter.ToHex6()

		assert.Equal(t, "#ab12ef", result)
	})
}

func TestConverter_ToHex6NoPrefix(t *testing.T) {
	t.Run("should_return_hex6_without_hash", func(t *testing.T) {
		converter, err := NewConverter("#ab12ef")
		require.NoError(t, err)

		result := converter.ToHex6NoPrefix()

		assert.Equal(t, "ab12ef", result)
	})
}

func TestConverter_ToHex0x(t *testing.T) {
	t.Run("should_return_0x_format", func(t *testing.T) {
		converter, err := NewConverter("#ab12ef")
		require.NoError(t, err)

		result := converter.ToHex0x()

		assert.Equal(t, "0xab12ef", result)
	})
}

func TestConverter_ToRGBNoPrefix(t *testing.T) {
	t.Run("should_return_rgb_format_without_prefix", func(t *testing.T) {
		converter, err := NewConverter("#ab12ef")
		require.NoError(t, err)

		result := converter.ToRGBNoPrefix()

		assert.Equal(t, "ab12ef", result)
	})
}

func TestConverter_ToFormat(t *testing.T) {
	t.Run("should_convert_to_hex6_format", func(t *testing.T) {
		converter, err := NewConverter("#ab12ef")
		require.NoError(t, err)

		result := converter.ToFormat("hex6")

		assert.Equal(t, "#ab12ef", result)
	})

	t.Run("should_convert_to_hex6_no_prefix_format", func(t *testing.T) {
		converter, err := NewConverter("#ab12ef")
		require.NoError(t, err)

		result := converter.ToFormat("hex6_no_prefix")

		assert.Equal(t, "ab12ef", result)
	})

	t.Run("should_convert_to_hex_0x_format", func(t *testing.T) {
		converter, err := NewConverter("#ab12ef")
		require.NoError(t, err)

		result := converter.ToFormat("hex_0x")

		assert.Equal(t, "0xab12ef", result)
	})

	t.Run("should_convert_to_rgb_no_prefix_format", func(t *testing.T) {
		converter, err := NewConverter("#ab12ef")
		require.NoError(t, err)

		result := converter.ToFormat("rgb_no_prefix")

		assert.Equal(t, "ab12ef", result)
	})

	t.Run("should_return_hex6_for_unknown_format", func(t *testing.T) {
		converter, err := NewConverter("#ab12ef")
		require.NoError(t, err)

		result := converter.ToFormat("unknown_format")

		assert.Equal(t, "#ab12ef", result)
	})
}

func TestConverter_FormatConsistency(t *testing.T) {
	t.Run("should_maintain_consistency_across_all_formats", func(t *testing.T) {
		converter, err := NewConverter("#ab12ef")
		require.NoError(t, err)

		assert.Equal(t, "#ab12ef", converter.ToHex6())
		assert.Equal(t, "ab12ef", converter.ToHex6NoPrefix())
		assert.Equal(t, "0xab12ef", converter.ToHex0x())
		assert.Equal(t, "ab12ef", converter.ToRGBNoPrefix())
	})

	t.Run("should_handle_edge_case_colors_consistently", func(t *testing.T) {
		converter, err := NewConverter("000000")
		require.NoError(t, err)

		assert.Equal(t, "#000000", converter.ToHex6())
		assert.Equal(t, "000000", converter.ToHex6NoPrefix())
		assert.Equal(t, "0x000000", converter.ToHex0x())
		assert.Equal(t, "000000", converter.ToRGBNoPrefix())
	})
}

func TestConverter_NilSafety(t *testing.T) {
	t.Run("should_handle_nil_converter_gracefully", func(t *testing.T) {
		var converter *Converter
		
		// These should not panic, but return safe defaults or handle gracefully
		assert.Equal(t, "", converter.ToHex6())
		assert.Equal(t, "", converter.ToHex6NoPrefix())
		assert.Equal(t, "", converter.ToHex0x())
		assert.Equal(t, "", converter.ToRGBNoPrefix())
		assert.Equal(t, "", converter.ToFormat("hex6"))
	})
}

func TestConverter_OutputValidation(t *testing.T) {
	t.Run("should_ensure_hex6_always_has_hash_prefix", func(t *testing.T) {
		converter, err := NewConverter("ffffff")
		require.NoError(t, err)
		
		result := converter.ToHex6()
		assert.True(t, strings.HasPrefix(result, "#"))
		assert.Len(t, result, 7) // # + 6 characters
	})

	t.Run("should_ensure_hex6_no_prefix_never_has_hash", func(t *testing.T) {
		converter, err := NewConverter("#ffffff")
		require.NoError(t, err)
		
		result := converter.ToHex6NoPrefix()
		assert.False(t, strings.HasPrefix(result, "#"))
		assert.Len(t, result, 6) // exactly 6 characters
	})

	t.Run("should_ensure_hex_0x_always_has_0x_prefix", func(t *testing.T) {
		converter, err := NewConverter("ffffff")
		require.NoError(t, err)
		
		result := converter.ToHex0x()
		assert.True(t, strings.HasPrefix(result, "0x"))
		assert.Len(t, result, 8) // 0x + 6 characters
	})

	t.Run("should_ensure_rgb_no_prefix_never_has_prefix", func(t *testing.T) {
		converter, err := NewConverter("#ffffff")
		require.NoError(t, err)
		
		result := converter.ToRGBNoPrefix()
		assert.False(t, strings.HasPrefix(result, "#"))
		assert.False(t, strings.HasPrefix(result, "0x"))
		assert.Len(t, result, 6) // exactly 6 characters
	})
}

func TestConverter_CaseConsistency(t *testing.T) {
	t.Run("should_always_output_lowercase", func(t *testing.T) {
		converter, err := NewConverter("#ABCDEF")
		require.NoError(t, err)
		
		assert.Equal(t, "#abcdef", converter.ToHex6())
		assert.Equal(t, "abcdef", converter.ToHex6NoPrefix())
		assert.Equal(t, "0xabcdef", converter.ToHex0x())
		assert.Equal(t, "abcdef", converter.ToRGBNoPrefix())
	})
}

func TestConverter_InputNormalization(t *testing.T) {
	t.Run("should_handle_input_with_whitespace", func(t *testing.T) {
		converter, err := NewConverter("  #ffffff  ")
		require.NoError(t, err)
		
		assert.Equal(t, "#ffffff", converter.ToHex6())
	})

	t.Run("should_handle_mixed_case_0x_prefix", func(t *testing.T) {
		converter, err := NewConverter("0XaBcDeF")
		require.NoError(t, err)
		
		assert.Equal(t, "#abcdef", converter.ToHex6())
	})
}

func TestConverter_Hex3Expansion(t *testing.T) {
	t.Run("should_expand_hex3_with_mixed_digits", func(t *testing.T) {
		converter, err := NewConverter("#a0f")
		require.NoError(t, err)
		
		assert.Equal(t, "#aa00ff", converter.ToHex6())
	})

	t.Run("should_expand_hex3_with_same_digits", func(t *testing.T) {
		converter, err := NewConverter("333")
		require.NoError(t, err)
		
		assert.Equal(t, "#333333", converter.ToHex6())
	})
}

func TestConverter_InvalidFormats(t *testing.T) {
	t.Run("should_return_error_for_empty_0x", func(t *testing.T) {
		converter, err := NewConverter("0x")
		
		assert.Error(t, err)
		assert.Nil(t, converter)
	})

	t.Run("should_return_error_for_0x_with_invalid_chars", func(t *testing.T) {
		converter, err := NewConverter("0xGGGGGG")
		
		assert.Error(t, err)
		assert.Nil(t, converter)
	})

	t.Run("should_return_error_for_0x_wrong_length", func(t *testing.T) {
		converter, err := NewConverter("0x12345")
		
		assert.Error(t, err)
		assert.Nil(t, converter)
	})
}

func TestConverter_ErrorMessages(t *testing.T) {
	t.Run("should_provide_clear_error_for_invalid_characters", func(t *testing.T) {
		_, err := NewConverter("#xyz123")
		
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid")
	})

	t.Run("should_provide_clear_error_for_wrong_length", func(t *testing.T) {
		_, err := NewConverter("#12345")
		
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "length")
	})
}

func TestConverter_NewMethods(t *testing.T) {
	t.Run("should_convert_to_rgba_string", func(t *testing.T) {
		converter, err := NewConverter("#ff0080")
		require.NoError(t, err)
		
		result := converter.ToRGBAString(0.8)
		assert.Equal(t, "rgba(255,0,128,0.8)", result)
	})
	
	t.Run("should_brighten_color", func(t *testing.T) {
		converter, err := NewConverter("#808080")
		require.NoError(t, err)
		
		result := converter.Brighten(0.3)
		assert.Equal(t, "#a6a6a6", result)
	})
	
	t.Run("should_mix_colors", func(t *testing.T) {
		conv1, err1 := NewConverter("#000000")
		conv2, err2 := NewConverter("#ffffff")
		require.NoError(t, err1)
		require.NoError(t, err2)
		
		result := conv1.Mix(conv2, 0.5)
		assert.Equal(t, "#7f7f7f", result)
	})
	
	t.Run("should_handle_nil_safely", func(t *testing.T) {
		var nilConverter *Converter
		
		assert.Equal(t, "", nilConverter.ToRGBAString(0.5))
		assert.Equal(t, "", nilConverter.Brighten(0.3))
		assert.Equal(t, "", nilConverter.Mix(nil, 0.5))
	})
}
