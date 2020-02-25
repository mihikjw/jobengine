package models

//Queued is the status a job is at before it is picked-up for processing
const Queued string = "queued"

//Inprogress is the status a job is at once it has been picked-up for processing
const Inprogress string = "inprogress"

//Complete is the status a job is in once processing has succesfully finished
const Complete string = "complete"

//Failed is the status a job is in once processing has unsuccesfully finished
const Failed string = "failed"
