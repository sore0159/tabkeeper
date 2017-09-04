### Tabkeeper

A simple web app for my wife and me to manage our shared finances.  This app serves a web page that allows us to add costs to a shared "tab" that are tracked and summed automatically.  Charges can be marked as "repeatable"; if so, when the charge is listed it will have a simple "repeat" button next to its listing.

Designed to be run on our home server with a reverse proxy allowing only home network access.  The UI, from the color coding to the page names, reflects this parochial focus.

Data is stored as json for human readability and error correcting.  For example, the system does not have any way to delete erroneous charges or segment the list: this is expected to be rare enough to simply perform via manual editing.
