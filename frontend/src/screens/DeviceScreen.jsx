import React, { useEffect, useState } from "react";
import {
  Button,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  Paper,
  Typography,
  Switch,
  CircularProgress,
  Dialog,
  DialogTitle,
  DialogContent,
  TextField,
  DialogActions,
  Box,
} from "@mui/material";
import {
  getDevices,
  getDeviceFunctions,
  sendDeviceCommand,
  mapSwitch,
  getMappedSwitches,
  updateMappedSwitch,
  deleteMappedSwitch,
  resetSwitchState,
} from "../services/api";
import { toast } from "react-toastify";

const DeviceList = () => {
  const [devices, setDevices] = useState([]);
  const [functions, setFunctions] = useState({});
  const [loading, setLoading] = useState(true);
  const [switchStates, setSwitchStates] = useState({}); // <== theo dõi trạng thái switch
  // ////////
  const [mapDialogOpen, setMapDialogOpen] = useState(false);
  const [currentMap, setCurrentMap] = useState({ deviceId: "", code: "" });
  const [mapInfo, setMapInfo] = useState({ name: "", ip: "" });
  const [mappedSwitches, setMappedSwitches] = useState([]);
  const [editDialogOpen, setEditDialogOpen] = useState(false);
  const [editMapInfo, setEditMapInfo] = useState({ id: "", name: "", ip: "" });

  const handleResetSwitchState = async () => {
    try {
      const res = await resetSwitchState();
      toast.success("Đã reset thành công!");
    } catch (err) {
      toast.error("Lỗi khi reset trạng thái");
    }
  };

  const handleOpenEditDialog = (item) => {
    setEditMapInfo({ id: item.ID || item.id, name: item.name, ip: item.ip });
    setEditDialogOpen(true);
  };

  const handleOpenMapDialog = (deviceId, code) => {
    setCurrentMap({ deviceId, code });
    setMapDialogOpen(true);
  };
  const handleUpdateMap = async () => {
    try {
      await updateMappedSwitch(editMapInfo.id, {
        name: editMapInfo.name,
        ip: editMapInfo.ip,
      });
      await fetchMappedSwitches();
      setEditDialogOpen(false);
    } catch (err) {
      console.error("Lỗi cập nhật map", err);
    }
  };

  const handleDeleteMap = async (id) => {
    if (window.confirm("Bạn có chắc muốn xoá?")) {
      try {
        await deleteMappedSwitch(id);
        await fetchMappedSwitches();
      } catch (err) {
        console.error("Xoá thất bại", err);
      }
    }
  };

  const handleSubmitMap = async () => {
    try {
      await mapSwitch({
        device_id: currentMap.deviceId,
        code: currentMap.code,
        name: mapInfo.name,
        ip: mapInfo.ip,
      });
      await fetchMappedSwitches(); // reload list
      setMapDialogOpen(false);
      setMapInfo({ name: "", ip: "" });
    } catch (err) {
      console.error("Map error", err);
    }
  };

  const fetchMappedSwitches = async () => {
    try {
      const res = await getMappedSwitches();
      setMappedSwitches(res.data);
    } catch (err) {
      console.error("Lỗi fetch mapped switches", err);
    }
  };

  useEffect(() => {
    fetchDevices();
    fetchMappedSwitches();
  }, []);

  const fetchDevices = async () => {
    try {
      const res = await getDevices();
      setDevices(res.data.result || []);
      setLoading(false);

      for (const dev of res.data.result) {
        const funcRes = await getDeviceFunctions(dev.id);
        const allFunctions = funcRes.data.result.properties;
        const switches = allFunctions.filter(
          (f) =>
            f.type === "bool" && !f.code.toLowerCase().includes("backlight")
        );

        setFunctions((prev) => ({
          ...prev,
          [dev.id]: switches,
        }));

        // Lưu trạng thái hiện tại của switch
        const newSwitchStates = {};
        switches.forEach((f) => {
          newSwitchStates[`${dev.id}_${f.code}`] = f.value || false;
        });
        setSwitchStates((prev) => ({ ...prev, ...newSwitchStates }));
      }
    } catch (err) {
      console.error(err);
      setLoading(false);
    }
  };

  const handleToggle = async (deviceId, code, value) => {
    try {
      await sendDeviceCommand(deviceId, { code, value });
      setSwitchStates((prev) => ({
        ...prev,
        [`${deviceId}_${code}`]: value,
      }));
    } catch (err) {
      console.error("❌ Failed to send command", err);
    }
  };

  const isMapped = (deviceId, code) => {
    return mappedSwitches.some(
      (item) => item.device_id === deviceId && item.code === code
    );
  };

  if (loading) return <CircularProgress />;

  return (
    <div style={{ width: "100%" }}>
      <Box
        display="flex"
        justifyContent="space-between"
        alignItems="center"
        mb={2}
      >
        <Typography variant="h4">Danh sách Devices</Typography>
        <Button
          variant="outlined"
          color="secondary"
          onClick={handleResetSwitchState}
        >
          Reset trạng thái switch
        </Button>
      </Box>

      <TableContainer component={Paper}>
        <Table>
          <TableHead>
            <TableRow>
              <TableCell>
                <strong>Device Name</strong>
              </TableCell>
              <TableCell>
                <strong>Device ID</strong>
              </TableCell>
              <TableCell>
                <strong>Model</strong>
              </TableCell>
              <TableCell>
                <strong>Switches</strong>
              </TableCell>
            </TableRow>
          </TableHead>
          <TableBody>
            {devices.map((device) => (
              <TableRow key={device.id}>
                <TableCell>{device.name}</TableCell>
                <TableCell>{device.id}</TableCell>
                <TableCell>{device.model}</TableCell>
                <TableCell>
                  {(functions[device.id] || []).map((f) => {
                    const key = `${device.id}_${f.code}`;
                    const isOn = switchStates[key] ?? false;

                    return (
                      <div key={f.code} style={{ marginBottom: 8 }}>
                        <Typography
                          variant="body2"
                          sx={{ display: "inline-block", mr: 2 }}
                        >
                          {f.name || f.code}
                        </Typography>
                        <Switch
                          checked={isOn}
                          onChange={(e) =>
                            handleToggle(device.id, f.code, e.target.checked)
                          }
                          color="primary"
                        />
                        <Button
                          variant={
                            isMapped(device.id, f.code)
                              ? "contained"
                              : "outlined"
                          }
                          color={
                            isMapped(device.id, f.code) ? "error" : "primary"
                          }
                          size="small"
                          onClick={() => handleOpenMapDialog(device.id, f.code)}
                          disabled={isMapped(device.id, f.code)}
                          sx={{ ml: 1 }}
                        >
                          {isMapped(device.id, f.code) ? "Đã Map" : "Map"}
                        </Button>
                      </div>
                    );
                  })}
                </TableCell>
              </TableRow>
            ))}
          </TableBody>
        </Table>
      </TableContainer>

      <Typography variant="h5" sx={{ mt: 4 }}>
        Danh sách Switch đã Map
      </Typography>
      <TableContainer component={Paper} sx={{ mt: 1 }}>
        <Table>
          <TableHead>
            <TableRow>
              <TableCell>
                <strong>Name</strong>
              </TableCell>
              <TableCell>
                <strong>Device ID</strong>
              </TableCell>
              <TableCell>
                <strong>Code</strong>
              </TableCell>
              <TableCell>
                <strong>IP</strong>
              </TableCell>
              <TableCell>
                <strong>Switches</strong>
              </TableCell>
            </TableRow>
          </TableHead>
          <TableBody>
            {mappedSwitches.map((item, idx) => (
              <TableRow key={idx}>
                <TableCell>{item.name}</TableCell>
                <TableCell>{item.device_id}</TableCell>
                <TableCell>{item.code}</TableCell>

                <TableCell>{item.ip}</TableCell>
                <TableCell>
                  <Switch
                    checked={
                      switchStates[`${item.device_id}_${item.code}`] ?? false
                    }
                    onChange={(e) =>
                      handleToggle(item.device_id, item.code, e.target.checked)
                    }
                    color="primary"
                  />
                  <Button
                    variant="outlined"
                    size="small"
                    onClick={() => handleOpenEditDialog(item)}
                    sx={{ ml: 1 }}
                  >
                    Sửa
                  </Button>
                  <Button
                    variant="outlined"
                    color="error"
                    size="small"
                    onClick={() => handleDeleteMap(item.ID || item.id)}
                    sx={{ ml: 1 }}
                  >
                    Xoá
                  </Button>
                </TableCell>
              </TableRow>
            ))}
          </TableBody>
        </Table>
      </TableContainer>

      <Dialog open={mapDialogOpen} onClose={() => setMapDialogOpen(false)}>
        <DialogTitle>Map Switch</DialogTitle>
        <DialogContent>
          <TextField
            margin="dense"
            label="Tên"
            fullWidth
            value={mapInfo.name}
            onChange={(e) => setMapInfo({ ...mapInfo, name: e.target.value })}
          />
          <TextField
            margin="dense"
            label="IP Address"
            fullWidth
            value={mapInfo.ip}
            onChange={(e) => setMapInfo({ ...mapInfo, ip: e.target.value })}
          />
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setMapDialogOpen(false)}>Huỷ</Button>
          <Button variant="contained" onClick={handleSubmitMap}>
            Lưu
          </Button>
        </DialogActions>
      </Dialog>

      {/* popup edit map */}
      <Dialog open={editDialogOpen} onClose={() => setEditDialogOpen(false)}>
        <DialogTitle>Sửa thông tin Switch</DialogTitle>
        <DialogContent>
          <TextField
            margin="dense"
            label="Tên"
            fullWidth
            value={editMapInfo.name}
            onChange={(e) =>
              setEditMapInfo({ ...editMapInfo, name: e.target.value })
            }
          />
          <TextField
            margin="dense"
            label="IP Address"
            fullWidth
            value={editMapInfo.ip}
            onChange={(e) =>
              setEditMapInfo({ ...editMapInfo, ip: e.target.value })
            }
          />
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setEditDialogOpen(false)}>Huỷ</Button>
          <Button variant="contained" onClick={handleUpdateMap}>
            Cập nhật
          </Button>
        </DialogActions>
      </Dialog>
    </div>
  );
};

export default DeviceList;
