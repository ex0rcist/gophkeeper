package top

import (
	"gophkeeper/internal/keeper/tui"

	blobEdit "gophkeeper/internal/keeper/tui/screens/blob_edit"
	cardEdit "gophkeeper/internal/keeper/tui/screens/card_edit"
	credentialEdit "gophkeeper/internal/keeper/tui/screens/credential_edit"
	"gophkeeper/internal/keeper/tui/screens/login"
	"gophkeeper/internal/keeper/tui/screens/menu"
	remoteeopen "gophkeeper/internal/keeper/tui/screens/remote_open"
	secretType "gophkeeper/internal/keeper/tui/screens/secret_type"
	storageBrowse "gophkeeper/internal/keeper/tui/screens/storage_browse"
	storageCreate "gophkeeper/internal/keeper/tui/screens/storage_create"
	storageOpen "gophkeeper/internal/keeper/tui/screens/storage_open"
	textEdit "gophkeeper/internal/keeper/tui/screens/text_edit"
	"gophkeeper/internal/keeper/tui/screens/welcome"
)

// Screen constructors. Inject dependencies if any
func prepareMakers(deps ModelDependencies) map[tui.Screen]tui.ScreenMaker {
	return map[tui.Screen]tui.ScreenMaker{
		tui.WelcomeScreen:        &welcome.WelcomeScreen{},
		tui.MenuScreen:           &menu.MenuScreen{},
		tui.StorageCreateScreen:  &storageCreate.StorageCreateScreen{},
		tui.StorageOpenScreen:    &storageOpen.StorageOpenScreen{},
		tui.StorageBrowseScreen:  &storageBrowse.StorageBrowseScreen{},
		tui.SecretTypeScreen:     &secretType.SecretTypeScreen{},
		tui.CredentialEditScreen: &credentialEdit.CredentialEditScreen{},
		tui.TextEditScreen:       &textEdit.TextEditScreen{},
		tui.CardEditScreen:       &cardEdit.CardEditScreen{},
		tui.BlobEditScreen:       &blobEdit.BlobEditScreen{},
		tui.FilePickScreen:       &blobEdit.FilePickScreen{},
		tui.LoginScreen:          &login.LoginScreen{},
		tui.RemoteOpenScreen:     &remoteeopen.RemoteOpenScreenMaker{Client: deps.Client},
	}
}
