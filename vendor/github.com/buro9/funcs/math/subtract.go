package math

// Subtract takes two numbers (a and b) and substracts b from a.
//
// The type of b will be cast to the same type as a.
func Subtract(a interface{}, b interface{}) interface{} {
	switch av := a.(type) {
	case float32:
		switch bv := b.(type) {
		case float32:
			return av - bv
		case float64:
			return av - float32(bv)
		case int:
			return av - float32(bv)
		case int32:
			return av - float32(bv)
		case int64:
			return av - float32(bv)
		default:
			return nil
		}
	case float64:
		switch bv := b.(type) {
		case float32:
			return av - float64(bv)
		case float64:
			return av - bv
		case int:
			return av - float64(bv)
		case int32:
			return av - float64(bv)
		case int64:
			return av - float64(bv)
		default:
			return nil
		}
	case int:
		switch bv := b.(type) {
		case float32:
			return av - int(bv)
		case float64:
			return av - int(bv)
		case int:
			return av - bv
		case int32:
			return av - int(bv)
		case int64:
			return av - int(bv)
		default:
			return nil
		}
	case int32:
		switch bv := b.(type) {
		case float32:
			return av - int32(bv)
		case float64:
			return av - int32(bv)
		case int:
			return av - int32(bv)
		case int32:
			return av - bv
		case int64:
			return av - int32(bv)
		default:
			return nil
		}
	case int64:
		switch bv := b.(type) {
		case float32:
			return av - int64(bv)
		case float64:
			return av - int64(bv)
		case int:
			return av - int64(bv)
		case int32:
			return av - int64(bv)
		case int64:
			return av - bv
		default:
			return nil
		}
	default:
		return nil
	}
}
