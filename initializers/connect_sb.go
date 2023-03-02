package initializers

import (
	"os"

	supabasestorageuploader "github.com/adityarizkyramadhan/supabase-storage-uploader"
)

var SupabaseClient supabasestorageuploader.SupabaseClientService

func ConnectSupabase() {
	supClient := supabasestorageuploader.NewSupabaseClient(
		os.Getenv("PROJECT_URL"),
		os.Getenv("PROJECT_API_KEYS"),
		os.Getenv("STORAGE_NAME"),
		os.Getenv("STORAGE_FOLDER"),
	)
	SupabaseClient = supClient
}
