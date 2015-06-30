namespace go translator

exception TranslatorError {
	1: TranslatorErrorCode error_code,
	2: string error_message,
}

enum TranslatorErrorCode {
	UNKNOWN_ERROR = 0,
	LOCALE_NOT_FOUND = 1,
}

service Translator {
	bool ping(),
	map<string,string> get_dictionary(1:string locale) throws (1: TranslatorError error),
}