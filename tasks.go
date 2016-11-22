/*
tasks is set of background tasks
*/
package core

/*
Task interface for background running tasks
 */
type Task interface {

	// Run executes task
	Run(config Config) error
}
