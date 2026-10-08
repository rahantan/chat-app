package pkg

// func VerifyJWTToken(secretKey string, tokenStr string) (*dto.UserResponse, error) {
// 	claims := &dto.JwtPayLoad{}
// 	token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (any, error) {
// 		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
// 			return nil, errors.New("invalid signature algorithm")
// 		}
// 		return []byte(secretKey), nil
// 	})

// 	if err != nil {
// 		return nil, errors.New("token is expired")
// 	}

// 	if !token.Valid {
// 		return nil, errors.New("invalid parse token")
// 	}

// 	return &claims.UserResponse, nil
// }

// func CreateJWTToken(secretKey string, expireTime time.Duration, user *dto.UserResponse) (string, error) {
// 	start := time.Now()
// 	payload := dto.JwtPayLoad{
// 		RegisteredClaims: jwt.RegisteredClaims{
// 			ExpiresAt: jwt.NewNumericDate(start.Add(expireTime)),
// 		},
// 		UserResponse: *user,
// 	}
// 	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, payload).SignedString([]byte(secretKey))
// 	return token, err
// }
