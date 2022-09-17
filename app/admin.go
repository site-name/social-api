package app

import (
	"io"
	"net/http"
	"os"
	"runtime/debug"

	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/modules/config"
	"github.com/sitename/sitename/modules/slog"
)

func (s *Server) GetLogs(page, perPage int) ([]string, *model.AppError) {
	var lines []string

	if s.Cluster != nil && *s.Config().ClusterSettings.Enable {
		if info := s.Cluster.GetMyClusterInfo(); info != nil {
			lines = append(lines, "-----------------------------------------------------------------------------------------------------------")
			lines = append(lines, "-----------------------------------------------------------------------------------------------------------")
			lines = append(lines, info.Hostname)
			lines = append(lines, "-----------------------------------------------------------------------------------------------------------")
			lines = append(lines, "-----------------------------------------------------------------------------------------------------------")
		} else {
			slog.Error("Could not get cluster info")
		}
	}

	melines, err := s.GetLogsSkipSend(page, perPage)
	if err != nil {
		return nil, err
	}

	lines = append(lines, melines...)

	if s.Cluster != nil && *s.Config().ClusterSettings.Enable {
		clines, err := s.Cluster.GetLogs(page, perPage)
		if err != nil {
			return nil, err
		}

		lines = append(lines, clines...)
	}

	return lines, nil
}

func (a *App) GetLogs(page, perPage int) ([]string, *model.AppError) {
	return a.Srv().GetLogs(page, perPage)
}

func (s *Server) GetLogsSkipSend(page, perPage int) ([]string, *model.AppError) {
	var lines []string

	if *s.Config().LogSettings.EnableFile {
		s.Log.Flush()
		logFile := config.GetLogFileLocation(*s.Config().LogSettings.FileLocation)
		file, err := os.Open(logFile)
		if err != nil {
			return nil, model.NewAppError("getLogs", "api.admin.file_read_error", nil, err.Error(), http.StatusInternalServerError)
		}

		defer file.Close()

		var newLine = []byte{'\n'}
		var lineCount int
		const searchPos = -1
		b := make([]byte, 1)
		var endOffset int64 = 0

		// if the file exists and it's last byte is '\n' - skip it
		var stat os.FileInfo
		if stat, err = os.Stat(logFile); err == nil {
			if _, err = file.ReadAt(b, stat.Size()-1); err == nil && b[0] == newLine[0] {
				endOffset = -1
			}
		}
		lineEndPos, err := file.Seek(endOffset, io.SeekEnd)
		if err != nil {
			return nil, model.NewAppError("getLogs", "api.admin.file_read_error", nil, err.Error(), http.StatusInternalServerError)
		}
		for {
			pos, err := file.Seek(searchPos, io.SeekCurrent)
			if err != nil {
				return nil, model.NewAppError("getLogs", "api.admin.file_read_error", nil, err.Error(), http.StatusInternalServerError)
			}

			_, err = file.ReadAt(b, pos)
			if err != nil {
				return nil, model.NewAppError("getLogs", "api.admin.file_read_error", nil, err.Error(), http.StatusInternalServerError)
			}

			if b[0] == newLine[0] || pos == 0 {
				lineCount++
				if lineCount > page*perPage {
					line := make([]byte, lineEndPos-pos)
					_, err := file.ReadAt(line, pos)
					if err != nil {
						return nil, model.NewAppError("getLogs", "api.admin.file_read_error", nil, err.Error(), http.StatusInternalServerError)
					}
					lines = append(lines, string(line))
				}
				if pos == 0 {
					break
				}
				lineEndPos = pos
			}

			if len(lines) == perPage {
				break
			}
		}

		for i, j := 0, len(lines)-1; i < j; i, j = i+1, j-1 {
			lines[i], lines[j] = lines[j], lines[i]
		}
	} else {
		lines = append(lines, "")
	}

	return lines, nil
}

func (a *App) GetLogsSkipSend(page, perPage int) ([]string, *model.AppError) {
	return a.Srv().GetLogsSkipSend(page, perPage)
}

func (a *App) GetClusterStatus() []*model.ClusterInfo {
	infos := make([]*model.ClusterInfo, 0)

	if a.Cluster() != nil {
		infos = a.Cluster().GetClusterInfos()
	}

	return infos
}

func (s *Server) InvalidateAllCaches() *model.AppError {
	debug.FreeOSMemory()
	s.InvalidateAllCachesSkipSend()

	if s.Cluster != nil {
		msg := &model.ClusterMessage{
			Event:            model.ClusterEventInvalidateAllCaches,
			SendType:         model.ClusterSendReliable,
			WaitForAllToSend: true,
		}

		s.Cluster.SendClusterMessage(msg)
	}

	return nil
}

func (s *Server) InvalidateAllCachesSkipSend() {
	slog.Info("Purging all caches")
	s.AccountService().ClearAllUsersSessionCacheLocal()
	s.StatusCache.Purge()
	s.Store.User().ClearCaches()
	s.Store.FileInfo().ClearCaches()
	// s.Store.Webhook().ClearCaches()
	// s.Store.Post().ClearCaches()
}

// serverBusyStateChanged is called when a CLUSTER_EVENT_BUSY_STATE_CHANGED is received.
func (s *Server) serverBusyStateChanged(sbs *model.ServerBusyState) {
	s.Busy.ClusterEventChanged(sbs)
	if sbs.Busy {
		slog.Warn("server busy state activitated via cluster event - non-critical services disabled", slog.Int64("expires_sec", sbs.Expires))
	} else {
		slog.Info("server busy state cleared via cluster event - non-critical services enabled")
	}
}
