package privet

type Printer struct {
	Name     string
	Location string
}
type PrinterInfo struct {
	Version         string   `json:"version"`
	Name            string   `json:"name"`
	Description     string   `json:"description"`
	Url             string   `json:"url"`
	Type            []string `json:"type"`
	Id              string   `json:"id"`
	DeviceState     string   `json:"device_state"`
	ConnectionState string   `json:"connection_state"`
	Manufacturer    string   `json:"manufacturer"`
	Model           string   `json:"model"`
	SerialNumber    string   `json:serial_number`
	Firmware        string   `json:"firmware"`
	Uptime          int      `json:"uptime"`
	SetupUrl        string   `json:"setup_url"`
	SupportUrl      string   `json:"support_url"`
	UpdateUrl       string   `json:"update_url"`
	XPrivetToken    string   `json:"x-privet-token"`
	Api             []string `json:"api"`
}

/*
job_id	(optional) Print job id. May be omitted for simple printing case (see above). Must match the one returned by the printer.
user_name	(optional) Human readable user name. This is not definitive, and should only be used for print job annotations. If job is re-posted to the Cloud Print service this string should be attached to the Cloud Print job.
client_name	(optional) Name of the client application making this request. For display purposes only. If job is re-posted to the Cloud Print service this string should be attached to the Cloud Print job.
job_name	(optional) Name of the print job to be recorded. If job is re-posted to the Cloud Print service this string should be attached to the Cloud Print job.
offline	(optional) Could only be "offline=1". In this case printer should only try printing offline (no re-post to Cloud Print server).
*/
type SubmitDocParams struct {
	JobId      string `json:"job_id"`
	UserName   string `json:"user_name"`
	ClientName string `json:"client_name"`
	JobName    string `json:"job_name"`
	Offline    int    `json:"offline"`
}

/*
job_id	string	ID of the newly created print job (simple printing) or job_id specified in the request (advanced printing).
expires_in	int	Number of seconds this print job is valid.
job_type	string	Content-type of the submitted document.
job_size	int 64 bit	Size of the print data in bytes.
job_name	string	(optional) Same job name as in input (if any).
*/
type SubmitDocResp struct {
	JobId     string `json:"job_id"`
	ExpiresIn int    `json:"expires_in"`
	JobType   string `json:"job_type"`
	JobSize   int    `json:"job_size"`
	JobName   string `json:"job_name"`
}
