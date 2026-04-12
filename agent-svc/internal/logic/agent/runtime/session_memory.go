package runtime

import "encoding/json"

func finalizeTaskSessionState(previous TaskSessionState, current TaskSessionState, conversationNo string) TaskSessionState {
	current = copyTaskSession(current)
	current.ConversationNo = conversationNo

	var (
		previousVersion = previous.SessionVersion
		currentVersion  = current.SessionVersion
	)
	if currentVersion > previousVersion {
		return current
	}
	if taskSessionStateChanged(previous, current) {
		current.SessionVersion = previousVersion + 1
		if current.SessionVersion == 0 {
			current.SessionVersion = 1
		}
		return current
	}
	if previousVersion > 0 {
		current.SessionVersion = previousVersion
	}
	return current
}

func taskSessionStateChanged(previous TaskSessionState, current TaskSessionState) bool {
	previous = sanitizeTaskSessionForFingerprint(previous)
	current = sanitizeTaskSessionForFingerprint(current)

	previousJSON, _ := json.Marshal(previous)
	currentJSON, _ := json.Marshal(current)
	return string(previousJSON) != string(currentJSON)
}

func sanitizeTaskSessionForFingerprint(session TaskSessionState) TaskSessionState {
	session.ConversationNo = ""
	session.SessionVersion = 0
	return session
}
