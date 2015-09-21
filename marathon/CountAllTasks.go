package marathon

////////////////////////////////////////////////////////////////////////////////////////
// CountAllTasks
//
// This does nothing for Marathon because we just don't
// use this data. We do use it for ECS though so this function is provided for interface consistency.
//
func (m *MarathonScheduler) CountAllTasks() (int, int, error) {
	return 0, 0, nil
}
