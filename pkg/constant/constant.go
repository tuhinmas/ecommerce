package constant

type ContextKey string

const (
	FiberContext  ContextKey = "fiberCtx"
	HeaderContext ContextKey = "headerCtx"
)

const (
	Like = "like"
	Pass = "pass"
)

const (
	Verified = "verified"
	Premium  = "premium"
)

const (
	ErrorPaymentMethodNotFound = "Payment method not found"
	ErrorAmountNotEnough       = "Amount is not enough"
	ErrorAmountTooMuch         = "Amount is too much"
	ErrorAlreadySwiped         = "You have already swiped 10 people, please upgrade your account"
	ErrorAlreadySwipedPerson   = "You have already swiped this person"
	ErrorEmailAlreadyExists    = "Email already exists"
	ErrorPhoneAlreadyExists    = "Phone already exists"
	ErrorEmailNotFound         = "Email not found"
	ErrorPasswordWrong         = "Password is wrong"
)

const (
	PAGE                = "page"
	LIMIT               = "limit"
	SEARCH              = "search"
	DefaultPage         = 1
	DefaultLimitPerPage = 10
)

const (
	RESPONSE = "Response"
	QUERY    = "Query"
	CONSUMER = "Consumer"
	PRODUCER = "Producer"
)
