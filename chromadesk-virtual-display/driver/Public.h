// Copyright 2026 ChromaDesk Authors
// ChromaDesk Virtual Display Driver - IOCTL definitions and data structures.
// Based on RustDeskIddDriver (MS-PL license).

#pragma once

#include <minwindef.h>
#include <winioctl.h>
#include <guiddef.h>

// ---------------------------------------------------------------------------
// IOCTL codes
// Using FILE_DEVICE_UNKNOWN + custom function codes to avoid collisions
// with RustDesk's IOCTL_CHANGER_BASE codes.
// ---------------------------------------------------------------------------

#define IOCTL_CHROMADESK_VD_PLUG_IN \
    CTL_CODE(FILE_DEVICE_UNKNOWN, 0x0901, METHOD_BUFFERED, FILE_ANY_ACCESS)

#define IOCTL_CHROMADESK_VD_PLUG_OUT \
    CTL_CODE(FILE_DEVICE_UNKNOWN, 0x0902, METHOD_BUFFERED, FILE_ANY_ACCESS)

#define IOCTL_CHROMADESK_VD_UPDATE_MODE \
    CTL_CODE(FILE_DEVICE_UNKNOWN, 0x0903, METHOD_BUFFERED, FILE_ANY_ACCESS)

#define IOCTL_CHROMADESK_VD_QUERY \
    CTL_CODE(FILE_DEVICE_UNKNOWN, 0x0904, METHOD_BUFFERED, FILE_ANY_ACCESS)

// ---------------------------------------------------------------------------
// Error status codes (severity = error, customer bit set)
// ---------------------------------------------------------------------------

#define STATUS_ERROR_ADAPTER_NOT_INIT      ((NTSTATUS)((3u << 30) + 11))
#define STATUS_ERROR_MONITOR_EXISTS        ((NTSTATUS)((3u << 30) + 51))
#define STATUS_ERROR_MONITOR_NOT_EXISTS    ((NTSTATUS)((3u << 30) + 52))
#define STATUS_ERROR_MONITOR_INVALID_PARAM ((NTSTATUS)((3u << 30) + 53))
#define STATUS_ERROR_MONITOR_OOM           ((NTSTATUS)((3u << 30) + 54))
#define STATUS_ERROR_INDEX_OOR             ((NTSTATUS)((3u << 30) + 55))

// ---------------------------------------------------------------------------
// Maximum virtual monitors per adapter
// ---------------------------------------------------------------------------

#define CHROMADESK_VD_MAX_MONITORS 4

// ---------------------------------------------------------------------------
// EDID indices for built-in monitor descriptors
// ---------------------------------------------------------------------------

#define MONITOR_EDID_MOD_DELL_S2719DGF 0
#define MONITOR_EDID_MOD_LENOVO_Y27fA  1

// ---------------------------------------------------------------------------
// IOCTL input/output structures
// ---------------------------------------------------------------------------

// IOCTL_CHROMADESK_VD_PLUG_IN input
typedef struct _ChromaDeskPlugInParams {
    UINT ConnectorIndex;   // 0 ~ CHROMADESK_VD_MAX_MONITORS-1
    UINT MonitorEDID;      // EDID index (0 = Dell, 1 = Lenovo, >=2 = EDID-less)
    GUID ContainerId;      // Unique container ID (caller generates via CoCreateGuid)
} ChromaDeskPlugInParams, *PChromaDeskPlugInParams;

// IOCTL_CHROMADESK_VD_PLUG_OUT input
typedef struct _ChromaDeskPlugOutParams {
    UINT ConnectorIndex;
} ChromaDeskPlugOutParams, *PChromaDeskPlugOutParams;

// IOCTL_CHROMADESK_VD_UPDATE_MODE input
typedef struct _ChromaDeskMonitorMode {
    DWORD Width;
    DWORD Height;
    DWORD Sync;  // Refresh rate in Hz
} ChromaDeskMonitorMode;

typedef struct _ChromaDeskUpdateModeParams {
    UINT ConnectorIndex;
    UINT ModeCount;
    ChromaDeskMonitorMode Modes[1];  // Variable-length array
} ChromaDeskUpdateModeParams, *PChromaDeskUpdateModeParams;

// IOCTL_CHROMADESK_VD_QUERY output
typedef struct _ChromaDeskMonitorInfo {
    UINT  ConnectorIndex;
    BOOL  Active;
} ChromaDeskMonitorInfo;

typedef struct _ChromaDeskQueryResult {
    UINT ActiveCount;
    ChromaDeskMonitorInfo Monitors[CHROMADESK_VD_MAX_MONITORS];
} ChromaDeskQueryResult, *PChromaDeskQueryResult;

// ---------------------------------------------------------------------------
// Device interface / symbolic link
// ---------------------------------------------------------------------------

#define CHROMADESK_VD_SYMBOLIC_LINK_NAME L"\\Device\\ChromaDeskVirtualDisplay"

// ---------------------------------------------------------------------------
// Backward-compatible typedefs for driver internal code (Driver.h / Driver.cpp)
// The driver core uses these names in PlugInMonitor/PlugOutMonitor/UpdateMonitorModes.
// ---------------------------------------------------------------------------

typedef struct _CtlPlugIn {
    UINT ConnectorIndex;
    UINT MonitorEDID;
    GUID ContainerId;
} CtlPlugIn, *PCtlPlugIn;

typedef struct _CtlPlugOut {
    UINT ConnectorIndex;
} CtlPlugOut, *PCtlPlugOut;

typedef struct _CtlMonitorModes {
    UINT ConnectorIndex;
    UINT ModeCount;
    struct {
        DWORD Width;
        DWORD Height;
        DWORD Sync;
    } Modes[1];
} CtlMonitorModes, *PCtlMonitorModes;
