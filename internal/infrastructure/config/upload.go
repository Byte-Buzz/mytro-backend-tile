package config

import "errors"

type UploadConfig struct {
	MaxImageSize     uint64
	MaxVideoSize     uint64
	MaxVideoDuration uint32
	URLTTL           uint32
}

func (c *UploadConfig) LoadFromEnv() {
	val := getIntEnv("MAX_IMAGE_SIZE", 15*1024*1024)
	if val > 0 {
		c.MaxImageSize = uint64(val)
	}

	val = getIntEnv("MAX_VIDEO_SIZE", 150*1024*1024)
	if val > 0 {
		c.MaxVideoSize = uint64(val)
	}

	val = getIntEnv("MAX_VIDEO_DURATION", 180)
	if val > 0 {
		c.MaxVideoDuration = uint32(val)
	}

	val = getIntEnv("URL_TTL", 5*60)
	if val > 0 {
		c.URLTTL = uint32(val)
	}
}

func (c *UploadConfig) Validate() error {
	if c.MaxImageSize < 1 {
		return errors.New("invalid max image size")
	}

	if c.MaxVideoSize < 1 {
		return errors.New("invalid max video size")
	}

	if c.MaxVideoDuration < 1 {
		return errors.New("invalid max video duration")
	}

	if c.URLTTL < 1 {
		return errors.New("invalid url ttl")
	}

	return nil
}
