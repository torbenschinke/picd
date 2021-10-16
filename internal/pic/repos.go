package pic

// StillCameraRepository is part of the domain core. Must not be used by presentation layers.
type StillCameraRepository interface {
	FindAllStillCameras() ([]StillCamera, error)
	CapturePhoto(settings Settings, dst []byte) ([]byte, error)
}
