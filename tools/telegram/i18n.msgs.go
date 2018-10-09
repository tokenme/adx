package telegram

var il8n_msgs = map[string]map[string]string{
	"Sorry, submitted verification code is invalid": map[string]string{
		"vi-VN": "Xin lỗi, mã xác minh không tồn tại", // 越南语
	},
	"Sorry, we have some internal server bug :(": map[string]string{
		"vi-VN": "Xin lỗi, máy chủ đang gặp phải một vài sự cố :(",
	},
	"Sorry, submitted verification code is invalid or expired": map[string]string{
		"vi-VN": "Xin lỗi, mã xác minh không tồn tại hoặc đã hết hạn",
	},
	"Sorry, you must submit your code in group": map[string]string{
		"vi-VN": "Xin lỗi, bạn cần phải dán mã xác minh vào group",
	},
	"Sorry, maybe you didn’t submit your wallet address?": map[string]string{
		"vi-VN": "Xin lỗi, có thể bạn đã không nhập Huobi UID của bạn?",
	},
	"Sorry, you already submitted in this airdrop and could not submit again": map[string]string{
		"vi-VN": "Xin lỗi, bạn đã đăng ký airdrop này rồi và không thể làm lại lần nữa",
	},
	"Great! please wait for the airdrop transaction complete.": map[string]string{
		"vi-VN": "Yes! Xin vui lòng chờ giao dịch airdrop được thực hiện",
	},
	"Your airdrop is in the pool, please wait for the transaction complete": map[string]string{
		"vi-VN": "Airdrop của bạn đã nằm trong pool, xin vui vòng chờ đợi giao dịch được thực hiện",
	},
}

func il8n_trans(s string, lang string) string {
	if lang == "" || lang == "en" {
		return s
	}

	if v, ok := il8n_msgs[s]; ok {
		if ret, ok := v[lang]; ok {
			return ret
		}
	}

	return s
}
