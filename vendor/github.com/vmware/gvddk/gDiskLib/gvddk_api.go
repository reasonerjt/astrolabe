/*
Copyright (c) 2018-2019 the gvddk contributors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package gDiskLib

// #cgo CFLAGS: -g -Wall
// #cgo LDFLAGS: -L../lib/vmware-vix-disklib/lib64 -lvixDiskLib
// #include "gvddk_c.h"
import "C"
import (
	"fmt"
	"unsafe"
	//"container/list"
)

//export GoLogWarn
func GoLogWarn(buf *C.char) {
	fmt.Println(C.GoString(buf))
}

func Init(majorVersion uint32, minorVersion uint32, dir string) VddkError {
	libDir := C.CString(dir)
	defer C.free(unsafe.Pointer(libDir))
	result := C.Init(C.uint32(majorVersion), C.uint32(minorVersion), libDir)
	if result != 0 {
		return NewVddkError(uint64(result), fmt.Sprintf("Initialize failed. The error code is %d.", result))
	}
	return nil
}

func prepareConnectParams(appGlobal ConnectParams) (*C.VixDiskLibConnectParams, []*C.char) {
	// Trans string to CString
	vmxSpec :=C.CString(appGlobal.vmxSpec)
	serverName :=C.CString(appGlobal.serverName)
	thumbPrint :=C.CString(appGlobal.thumbPrint)
	userName :=C.CString(appGlobal.userName)
	password :=C.CString(appGlobal.password)
	fcdId :=C.CString(appGlobal.fcdId)
	ds :=C.CString(appGlobal.ds)
	fcdssId :=C.CString(appGlobal.fcdssId)
	cookie :=C.CString(appGlobal.cookie)
	var cParams = []*C.char{vmxSpec, serverName, thumbPrint, userName, password, fcdId, ds, fcdssId, cookie}
	// Construct connparams which can be c wrapper used directly
	var cnxParams *C.VixDiskLibConnectParams = C.VixDiskLib_AllocateConnectParams()
	if appGlobal.fcdId != "" {
		cnxParams.specType = C.VIXDISKLIB_SPEC_VSTORAGE_OBJECT
		C.Params_helper(cnxParams, fcdId, ds, fcdssId, true, false)
	} else if appGlobal.vmxSpec != "" {
		cnxParams.specType = C.VIXDISKLIB_SPEC_VMX
		cnxParams.vmxSpec = vmxSpec
	}
	cnxParams.thumbPrint = thumbPrint
	cnxParams.serverName = serverName
	if appGlobal.cookie == "" {
		cnxParams.credType = C.VIXDISKLIB_CRED_UID
		C.Params_helper(cnxParams, cookie, userName, password, false, false)
	} else {
		cnxParams.credType = C.VIXDISKLIB_CRED_SESSIONID
		C.Params_helper(cnxParams, cookie, userName, password, false, true)
	}
	return cnxParams, cParams
}

func freeParams(params []*C.char) {
	for i, _ := range(params) {
		C.free(unsafe.Pointer(params[i]))
	}
	return
}

func Connect(appGlobal ConnectParams) (VixDiskLibConnection, VddkError) {
	var connection VixDiskLibConnection
	cnxParams, toFree := prepareConnectParams(appGlobal)
	defer freeParams(toFree)
	err := C.Connect(cnxParams, &connection.conn)
	if err != 0 {
		return VixDiskLibConnection{}, NewVddkError(uint64(err), fmt.Sprintf("Connect failed. The error code is %d.", err))
	}

	return connection, nil
}

func ConnectEx(appGlobal ConnectParams) (VixDiskLibConnection, VddkError) {
	var connection VixDiskLibConnection
	cnxParams, toFree := prepareConnectParams(appGlobal)
	defer freeParams(toFree)
	modes := C.CString(appGlobal.mode)
	defer C.free(unsafe.Pointer(modes))
	err := C.ConnectEx(cnxParams, C._Bool(appGlobal.readOnly), modes, &connection.conn)
	if err != 0 {
		return VixDiskLibConnection{}, NewVddkError(uint64(err), fmt.Sprintf("ConnectEx failed. The error code is %d.", err))
	}

	return connection, nil
}

func PrepareForAccess(appGlobal ConnectParams) VddkError {
	name := C.CString(appGlobal.identity)
	defer C.free(unsafe.Pointer(name))
	cnxParams, toFree := prepareConnectParams(appGlobal)
	defer freeParams(toFree)
	result := C.PrepareForAccess(cnxParams, name)
	if result != 0 {
		return NewVddkError(uint64(result), fmt.Sprintf("Prepare for access failed. The error code is %d.", result))
	}
	return nil
}

func Open(conn VixDiskLibConnection, params ConnectParams) (VixDiskLibHandle, VddkError) {
	var dli VixDiskLibHandle
	filePath := C.CString(params.path)
	defer C.free(unsafe.Pointer(filePath))
	res := C.Open(conn.conn, filePath, C.uint32(params.flag))
	dli.dli = res.dli
	if res.err != 0 {
		return dli, NewVddkError(uint64(res.err), fmt.Sprintf("Open virtual disk file failed. The error code is %d.", res.err))
	}
	return dli, nil
}

func EndAccess(appGlobal ConnectParams) VddkError {
	name := C.CString(appGlobal.identity)
	defer C.free(unsafe.Pointer(name))
	cnxParams, toFree := prepareConnectParams(appGlobal)
	result := C.VixDiskLib_EndAccess(cnxParams, name)
	freeParams(toFree)
	if result != 0 {
		return NewVddkError(uint64(result), fmt.Sprintf("End access failed. The error code is %d.", result))
	}
	return nil
}

func Disconnect(connection VixDiskLibConnection) VddkError {
	res := C.VixDiskLib_Disconnect(connection.conn)
	if res != 0 {
		return NewVddkError(uint64(res), fmt.Sprintf("Disconnect failed. The error code is %d.", res))
	}
	return nil
}

//func Exit() {
//	C.VixDiskLib_Exit()
//}

func Attach(childHandle VixDiskLibHandle, parentHandle VixDiskLibHandle) VddkError {
	res := C.VixDiskLib_Attach(childHandle.dli, parentHandle.dli)
	if res != 0 {
		return NewVddkError(uint64(res), fmt.Sprintf("Attach child disk chain to the parent disk chain failed. The error code is %d.", res))
	}
	return nil
}

func CheckRepair(connection VixDiskLibConnection, filename string, repair bool) VddkError {
	file := C.CString(filename)
	defer C.free(unsafe.Pointer(file))
	res := C.CheckRepair(connection.conn, file, C._Bool(repair))
	if res != 0 {
		return NewVddkError(uint64(res), fmt.Sprintf("Check repair failed. The error code is %d.", res))
	}
	return nil
}

func Cleanup(appGlobal ConnectParams, numCleanUp uint32, numRemaining uint32) VddkError {
	cnxParams, toFree := prepareConnectParams(appGlobal)
	defer freeParams(toFree)
	res := C.Cleanup(cnxParams, C.uint32(numCleanUp), C.uint32(numRemaining))
	if res != 0 {
		return NewVddkError(uint64(res), fmt.Sprintf("Clean up failed. The error code is %d.", res))
	}
	return nil
}

//func Clone(dstConnection VixDiskLibConnection, dstPath string, srcConnection VixDiskLibConnection, srcPath string,
//	       params VixDiskLibCreateParams, progressCallbackData string, overWrite bool) VddkError {
//	dst := C.CString(dstPath)
//	defer C.free(unsafe.Pointer(dst))
//	src := C.CString(srcPath)
//	defer C.free(unsafe.Pointer(src))
//	createParams := prepareCreateParams(params)
//	cstr := C.CString(progressCallbackData)
//	defer C.free(unsafe.Pointer(cstr))
//	res := C.Clone(dstConnection.conn, dst, srcConnection.conn, src, createParams, cstr, C._Bool(overWrite))
//	if res != 0 {
//		return NewVddkError(uint64(res), fmt.Sprintf("Clone a virtual disk failed. The error code is %d.", res))
//	}
//	return nil
//}

func prepareCreateParams(createSpec VixDiskLibCreateParams) *C.VixDiskLibCreateParams {
	var createParams *C.VixDiskLibCreateParams
	createParams.diskType = C.VixDiskLibDiskType(createSpec.diskType)
	createParams.adapterType = C.VixDiskLibAdapterType(createSpec.adapterType)
	createParams.hwVersion = C.uint16(createSpec.hwVersion)
	createParams.capacity = C.VixDiskLibSectorType(createSpec.capacity)
	return createParams
}

//func Create(connection VixDiskLibConnection, path string, createParams VixDiskLibCreateParams, progressCallbackData string) VddkError {
//	pathName := C.CString(path)
//	defer C.free(unsafe.Pointer(pathName))
//	createSpec := prepareCreateParams(createParams)
//	cstr := C.CString(progressCallbackData)
//	defer C.free(unsafe.Pointer(cstr))
//	res := C.Create(connection.conn, pathName, createSpec, cstr)
//	if res != 0 {
//		return NewVddkError(uint64(res), fmt.Sprintf("Create a virtual disk failed. The error code is %d.", res))
//	}
//	return nil
//}

//func CreateChild(diskHandle DiskHandle, childPath string, diskType VixDiskLibDiskType, progressCallbackData string) VddkError {
//	child := C.CString(childPath)
//	defer C.free(unsafe.Pointer(child))
//	cstr := C.CString(progressCallbackData)
//	defer C.free(unsafe.Pointer(cstr))
//	res := C.CreateChild(diskHandle.dli, child, C.VixDiskLibDiskType(diskType), cstr)
//	if res != 0 {
//		return NewVddkError(uint64(res), fmt.Sprintf("Create child virtual disk failed. The error code is %d.", res))
//	}
//	return nil
//}

func FreeErrorText(vixErrorMsg string) {
	errorMsg := C.CString(vixErrorMsg)
	defer C.free(unsafe.Pointer(errorMsg))
	C.VixDiskLib_FreeErrorText(errorMsg)
}

func FreeInfo(diskInfo *VixDiskLibInfo) {
	dliInfo, toFree := createDiskInfo(diskInfo)
	defer freeParams(toFree)
	C.VixDiskLib_FreeInfo(dliInfo)
}

func GetErrorText(error VddkError, locale string) string {
	lc := C.CString(locale)
	defer C.free(unsafe.Pointer(lc))
	res := C.VixDiskLib_GetErrorText(C.VixError(error.VixErrorCode()), lc)
	return C.GoString(res)
}

func createDiskInfo(diskInfo *VixDiskLibInfo) (*C.VixDiskLibInfo, []*C.char) {
	var dliInfo *C.VixDiskLibInfo
	var bios C.VixDiskLibGeometry
	var phys C.VixDiskLibGeometry
	bios.cylinders = C.uint32(diskInfo.biosGeo.cylinders)
	bios.heads = C.uint32(diskInfo.biosGeo.heads)
	bios.sectors = C.uint32(diskInfo.biosGeo.sectors)
	phys.cylinders = C.uint32(diskInfo.physGeo.cylinders)
	phys.heads = C.uint32(diskInfo.physGeo.heads)
	phys.sectors = C.uint32(diskInfo.physGeo.sectors)
	dliInfo.biosGeo = bios
	dliInfo.physGeo = phys
	dliInfo.capacity = C.VixDiskLibSectorType(diskInfo.capacity)
	dliInfo.adapterType = C.VixDiskLibAdapterType(diskInfo.adapterType)
	dliInfo.numLinks = C.int(diskInfo.numLinks)
	dliInfo.parentFileNameHint = C.CString(diskInfo.parentFileNameHint)
	dliInfo.uuid = C.CString(diskInfo.uuid)
	var cParams = []*C.char{dliInfo.parentFileNameHint, dliInfo.uuid}
	return dliInfo, cParams
}

func GetInfo(handle VixDiskLibHandle, diskInfo *VixDiskLibInfo) VddkError {
	dliInfo, toFree := createDiskInfo(diskInfo)
	defer freeParams(toFree)
	res := C.GetInfo(handle.dli, dliInfo)
	if res != 0 {
		return NewVddkError(uint64(res), fmt.Sprintf("GetInfo failed. The error code is %d.", res))
	}
	return nil
}

//func Grow(connection VixDiskLibConnection, path string, capacity VixDiskLibSectorType, updateGeometry bool, callbackData string) VddkError {
//	filePath := C.CString(path)
//	defer C.free(unsafe.Pointer(filePath))
//	cstr := C.CString(callbackData)
//	defer C.free(unsafe.Pointer(cstr))
//	res := C.Grow(connection.conn, filePath, C.VixDiskLibSectorType(capacity), C._Bool(updateGeometry), cstr)
//	if res != 0 {
//		return NewVddkError(uint64(res), fmt.Sprintf("Grow failed. The error code is %d.", res))
//	}
//	return nil
//}

func ListTransportModes() string {
	res := C.VixDiskLib_ListTransportModes()
	modes := C.GoString(res)
	return modes
}

func Rename(srcFileName string, dstFileName string) VddkError {
	src := C.CString(srcFileName)
	defer C.free(unsafe.Pointer(src))
	dst := C.CString(dstFileName)
	defer C.free(unsafe.Pointer(dst))
	res := C.VixDiskLib_Rename(src, dst)
	if res != 0 {
		return NewVddkError(uint64(res), fmt.Sprintf("Rename failed. The error code is %d.", res))
	}
	return nil
}

func SpaceNeededForClone(srcHandle VixDiskLibHandle, diskType VixDiskLibDiskType, spaceNeeded uint64) VddkError {
	space := C.uint64(spaceNeeded)
	res := C.VixDiskLib_SpaceNeededForClone(srcHandle.dli, C.VixDiskLibDiskType(diskType), &space)
	if res != 0 {
		return NewVddkError(uint64(res), fmt.Sprintf("Get space needed for clone failed. The error code is %d.", res))
	}
	return nil
}

func Unlink(connection VixDiskLibConnection, path string) VddkError {
	delete := C.CString(path)
	defer C.free(unsafe.Pointer(delete))
	res := C.VixDiskLib_Unlink(connection.conn, delete)
	if res != 0 {
		return NewVddkError(uint64(res), fmt.Sprintf("Delete the virtual disk including all the extents failed. The error code is %d.", res))
	}
	return nil
}

func Shrink(diskHandle VixDiskLibHandle, progressCallbackData string) VddkError {
	cstr := C.CString(progressCallbackData)
	defer C.free(unsafe.Pointer(cstr))
	res := C.Shrink(diskHandle.dli, unsafe.Pointer(&cstr))
	if res != 0 {
		return NewVddkError(uint64(res), fmt.Sprintf("Shrink failed. The error code is %d.", res))
	}
	return nil
}

func Defragment(diskHandle VixDiskLibHandle, progressCallbackData string) VddkError {
	cstr := C.CString(progressCallbackData)
	defer C.free(unsafe.Pointer(cstr))
	res := C.Defragment(diskHandle.dli, unsafe.Pointer(&cstr))
	if res != 0 {
		return NewVddkError(uint64(res), fmt.Sprintf("Defragment failed. The error code is %d.", res))
	}
	return nil
}

func GetTransportMode(diskHandle VixDiskLibHandle) string {
	res := C.VixDiskLib_GetTransportMode(diskHandle.dli)
	mode := C.GoString(res)
	return mode
}

func GetMetadataKeys(diskHandle VixDiskLibHandle, buf []byte, bufLen uint, requireLen uint) VddkError {
	cbuf := ((*C.char)(unsafe.Pointer(&buf[0])))
	res := C.GetMetadataKeys(diskHandle.dli, cbuf, C.size_t(bufLen), C.size_t(requireLen))
	if res != 0 {
		return NewVddkError(uint64(res), fmt.Sprintf("GetMetadataKeys failed. The error code is %d.", res))
	}
	return nil
}

func Close(diskHandle VixDiskLibHandle) VddkError {
	res := C.VixDiskLib_Close(diskHandle.dli)
	if res != 0 {
		return NewVddkError(uint64(res), fmt.Sprintf("Close virtual disk failed. The error code is %d.", res))
	}
	return nil
}

func WriteMetadata(diskHandle VixDiskLibHandle, key string, val string) VddkError {
	w_key := C.CString(key)
	defer C.free(unsafe.Pointer(w_key))
	w_val := C.CString(val)
	defer C.free(unsafe.Pointer(w_val))
	res := C.VixDiskLib_WriteMetadata(diskHandle.dli, w_key, w_val)
	if res != 0 {
		return NewVddkError(uint64(res), fmt.Sprintf("Write meta data failed. The error code is %d.", res))
	}
	return nil
}

func ReadMetadata(diskHandle VixDiskLibHandle, key string, buf []byte, bufLen uint, requiredLen uint) VddkError {
	readKey := C.CString(key)
	defer C.free(unsafe.Pointer(readKey))
	cbuf := ((*C.char)(unsafe.Pointer(&buf[0])))
	required := C.size_t(requiredLen)
	res := C.VixDiskLib_ReadMetadata(diskHandle.dli, readKey, cbuf, C.size_t(bufLen), &required)
	if res != 0 {
		return NewVddkError(uint64(res), fmt.Sprintf("Read meta data from virtual disk file failed. The error code is %d.", res))
	}
	return nil
}

func Read(diskHandle VixDiskLibHandle, startSector uint64, numSectors uint64, buf []byte) VddkError {
	cbuf := ((*C.uint8)(unsafe.Pointer(&buf[0])))
	res := C.VixDiskLib_Read(diskHandle.dli, C.VixDiskLibSectorType(startSector), C.VixDiskLibSectorType(numSectors), cbuf)
	if res != 0 {
		return NewVddkError(uint64(res), fmt.Sprintf("Read from virtual disk file failed. The error code is %d.", res))
	}
	return nil
}

func Write(diskHandle VixDiskLibHandle, startSector uint64, numSectors uint64, buf []byte) VddkError {
	cbuf := ((*C.uint8)(unsafe.Pointer(&buf[0])))
	res := C.VixDiskLib_Write(diskHandle.dli, C.VixDiskLibSectorType(startSector), C.VixDiskLibSectorType(numSectors), cbuf)
	if res != 0 {
		return NewVddkError(uint64(res), fmt.Sprintf("Write to virtual disk file failed. The error code is %d.", res))
	}
	return nil
}
