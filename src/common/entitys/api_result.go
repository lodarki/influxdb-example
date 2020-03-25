package entitys

type ApiResult struct {
	Code   ErrorCode   `json:"code"`
	Msg    string      `json:"msg"`
	Result interface{} `json:"result"`
}

type ErrorCode int

const (
	ERROR_CODE_SUCCESS                        ErrorCode = 0
	ERROR_CODE_GENERATE_HASHED_PASSWORD       ErrorCode = 1
	ERROR_CODE_GENERATE_TOKEN                 ErrorCode = 2
	ERROR_CODE_REMOVE_CAPTCHA                 ErrorCode = 3
	ERROR_CODE_UNKNOWN                        ErrorCode = 4
	ERROR_CODE_INTERNAL                       ErrorCode = 5
	ERROR_CODE_DB                             ErrorCode = 6
	ERROR_CODE_WEB_SOCKET_REGISTER_FAILED     ErrorCode = 7
	ERROR_CODE_WEB_SOCKET_SEND_MSG_FAILED     ErrorCode = 8
	ERROR_CODE_TIMESTAMP_NO_SYNC              ErrorCode = 9
	ERROR_CODE_DECODE_FAILED                  ErrorCode = 10
	ERROR_CODE_SEND_SMS_FAILED                ErrorCode = 11
	ERROR_CODE_PROCESS_TIMEOUT                ErrorCode = 12
	ERROR_CODE_COMMON_UNLOGIN                 ErrorCode = 201
	ERROR_CODE_PERMISSION_DENIED              ErrorCode = 202
	ERROR_CODE_SIGN_ERROR                     ErrorCode = 203
	ERROR_CODE_PARAMS_ERROR                   ErrorCode = 204
	ERROR_CODE_NOT_DATA_FOUND                 ErrorCode = 205
	ERROR_CODE_UNSUPPORTED_REQUEST            ErrorCode = 400
	ERROR_CODE_UNSUPPORTED_PROTOCOL           ErrorCode = 401
	ERROR_CODE_OPERATE_FAILED                 ErrorCode = 402
	ERROR_CODE_INVALID_ARGUMENT               ErrorCode = 403
	ERROR_CODE_USER_EXIST                     ErrorCode = 600
	ERROR_CODE_USER_NOT_EXIST                 ErrorCode = 601
	ERROR_CODE_USER_PASSWORD_NO_MATCH         ErrorCode = 602
	ERROR_CODE_INVALID_ACCESS_TOKEN           ErrorCode = 603
	ERROR_CODE_INVALID_CLIENT_TYPE            ErrorCode = 604
	ERROR_CODE_INVALID_LANG                   ErrorCode = 605
	ERROR_CODE_INVALID_MOBILE                 ErrorCode = 606
	ERROR_CODE_INVALID_CAPTCHA                ErrorCode = 607
	ERROR_CODE_SEND_CAPTCHA                   ErrorCode = 608
	ERROR_CODE_GET_CAPTCHA                    ErrorCode = 609
	ERROR_CODE_INVALID_PASSWORD               ErrorCode = 610
	ERROR_CODE_GENERATE_BASE64_CAPTCHA_FAILED ErrorCode = 611
	ERROR_CODE_VERIFY_BASE64_CAPTCHA_FAILED   ErrorCode = 612
	ERROR_CODE_REQUIRE_BASE64_CAPTCHA         ErrorCode = 613
	ERROR_CODE_SET_CAPTCHA                    ErrorCode = 614
	ERROR_CODE_NEED_UPGRADE_MONITOR           ErrorCode = 615
	ERROR_CODE_NO_PERMISSION_MONITOR_ID       ErrorCode = 616
	ERROR_CODE_MENU_NOT_EXIST                 ErrorCode = 700
	ERROR_CODE_MENU_ALREADY_EXIST             ErrorCode = 701
	ERROR_CODE_MENU_HAS_CHILD_NODE            ErrorCode = 702
	ERROR_CODE_MONITOR_NOT_EXIST              ErrorCode = 800
	ERROR_CODE_MONITOR_PERMISSION_DENIED      ErrorCode = 801
	ERROR_CODE_MONITOR_PARAMS_ERROR           ErrorCode = 802
	ERROR_CODE_MONITOR_BOUND                  ErrorCode = 803
	ERROR_CODE_MONITOR_MAC_DUPLICATE          ErrorCode = 804
	ERROR_CODE_MONITOR_NOT_BOUND              ErrorCode = 805
	ERROR_CODE_MONITOR_BACKUP_DIR_ERROR       ErrorCode = 806
	ERROR_CODE_MONITOR_BACKUP_FILE_SAVE_FAIL  ErrorCode = 807
	ERROR_CODE_MONITOR_BACKUP_FILE_NOT_FOUND  ErrorCode = 808
	ERROR_CODE_MONITOR_BACKUP_FILE_TOO_MANY   ErrorCode = 809
	ERROR_CODE_MONITOR_BACKUP_FILE_TYPE_ERROR ErrorCode = 810
	ERROR_CODE_MONITOR_BACKUP_FILE_MD5_ERROR  ErrorCode = 811
	ERROR_CODE_MONITOR_NOT_RESPONSE           ErrorCode = 812
	ERROR_CODE_MONITOR_SELECT                 ErrorCode = 813
	ERROR_CODE_MONITOR_VNC_CHANGE_PWD_FAILED  ErrorCode = 814
	ERROR_CODE_MONITOR_ALREADY_BIND           ErrorCode = 815
	ERROR_CODE_MONITOR_USER_ALREADY_BIND      ErrorCode = 816
	ERROR_CODE_MONITOR_VERSION_NOT_EXIST      ErrorCode = 820
	ERROR_CODE_NOT_HAVE_MONITOR               ErrorCode = 900
	ERROR_CODE_ROLE_ERROR                     ErrorCode = 901
	ERROR_CODE_ERROR_GRANT_PARAMS             ErrorCode = 902
	ERROR_CODE_GRANT_SELF                     ErrorCode = 903
	ERROR_CODE_WALLET_NAME_DUPLICATE          ErrorCode = 904
	ERROR_CODE_WALLET_ADDRESS_DUPLICATE       ErrorCode = 905
	ERROR_CODE_MONITOR_NO_FOUND               ErrorCode = 1000
	ERROR_CODE_MINER_IN_PROCESSING            ErrorCode = 1001
	ERROR_CODE_MONITOR_AUTHENTICATION_DENY    ErrorCode = 1002
	ERROR_CODE_TASK_NOT_BELONG_MONITOR        ErrorCode = 1003
	ERROR_CODE_TASK_FAIL                      ErrorCode = 1004
	ERROR_CODE_MONITOR_GROUP_NOT_FOUND        ErrorCode = 1005
	ERROR_CODE_MONITOR_GROUP_HAVE_MINER       ErrorCode = 1006
	ERROR_CODE_WALLET_NOT_ALLOW               ErrorCode = 1007
	ERROR_CODE_LOG_NO_FOUND                   ErrorCode = 1100
	ERROR_CODE_ROLE_HAS_BEEN_USERD            ErrorCode = 1101
	ERROR_CODE_ADMIN_USER_ALREADY_EXIST       ErrorCode = 1102
	ERROR_CODE_ADMIN_USER_IS_OFF              ErrorCode = 1103
)
