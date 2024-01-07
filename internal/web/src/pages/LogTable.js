import React, { useState, useEffect } from 'react';
import { Navbar } from '../components/Navbar';
import { DataGrid } from '@mui/x-data-grid';
import { Grid, TextField } from '@mui/material';
import { AdapterDayjs } from '@mui/x-date-pickers/AdapterDayjs';
import { LocalizationProvider } from '@mui/x-date-pickers/LocalizationProvider';
import { DateTimePicker } from '@mui/x-date-pickers/DateTimePicker';


const LogTable = () => {

    const [pageState, setPageState] = useState({
        isLoading: true,
        data: [],
        total: 0,
        page: 0,
        pageSize: 10,
    });

    const [filterState, setFilterState] = useState({
        level: '',
        message: '',
        resourceId: '',
        timestampStart: null,
        timestampEnd: null,
    });

    const columns = [
        { field: 'id', headerName: 'ID', width: 100 },
        { field: 'level', headerName: 'Level', width: 130 },
        { field: 'message', headerName: 'Message', width: 400 },
        { field: 'resourceId', headerName: 'Resource ID', width: 200 },
        { field: 'timestamp', headerName: 'Timestamp', width: 200 },
        { field: 'traceId', headerName: 'Trace ID', width: 200 },
        { field: 'spanId', headerName: 'Span ID', width: 200 },
        { field: 'commit', headerName: 'Commit', width: 200 },
        { field: 'parentResourceId', headerName: 'Parent Resource ID', width: 200 },
    ];
   
    useEffect(() => {
        const fetchLogs = async () => {
            setPageState({
                ...pageState,
                isLoading: true,
            });

            const params = {
                page: pageState.page,
                limit: pageState.pageSize,
                level: filterState.level,
                message: filterState.message,
                resourceId: filterState.resourceId,
                timestampStart: filterState.timestampStart === null ? '' : filterState.timestampStart,
                timestampEnd: filterState.timestampEnd === null ? '' : filterState.timestampEnd,
            };

            const res = await fetch(`http://localhost:1323/internal/logs?` + new URLSearchParams(params).toString());
            const data = await res.json();
            setPageState({
                ...pageState,
                isLoading: false,
                data: data.logs,
                total: data.total,
            });
        }
        fetchLogs();

    }, [pageState.page, pageState.pageSize, filterState]);




    return (
        <div>
            <Navbar
                title="Logs Visualization"
            />

            {/* Filter items */}
            <div style={{ display: 'flex', justifyContent: 'center', alignItems: 'center' }}>
                <Grid container spacing={2} style={{ padding: '20px', maxWidth: '1000px' }}>
                    <Grid item xs={12} md={4}>
                        <TextField
                            id="outlined-basic"
                            label="Search Level"
                            variant="outlined"
                            value={filterState.level}
                            onChange={(e) => {
                                setFilterState({
                                    ...filterState,
                                    level: e.target.value,
                                });
                            }}
                        />
                    </Grid>
                    <Grid item xs={12} md={4}>
                        <TextField
                            id="outlined-basic"
                            label="Search Message"
                            variant="outlined"
                            value={filterState.message}
                            onChange={(e) => {
                                setFilterState({
                                    ...filterState,
                                    message: e.target.value,
                                });
                            }}
                        />
                    </Grid>
                    <Grid item xs={12} md={4}>
                        <TextField
                            id="outlined-basic"
                            label="Search Resource ID"
                            variant="outlined"
                            value={filterState.resourceId}
                            onChange={(e) => {
                                setFilterState({
                                    ...filterState,
                                    resourceId: e.target.value,
                                });
                            }}
                        />
                    </Grid>
                </Grid>
            </div>

            <div style={{ display: 'flex', justifyContent: 'center', alignItems: 'center' }}>
                <Grid container  style={{ padding: '20px', maxWidth: '1000px' }}>
                    {/* TimePicker for Timestamp filter */}
                    <LocalizationProvider dateAdapter={AdapterDayjs}>
                        <Grid item xs={12} md={4}>
                            <DateTimePicker
                                label="Start Timestamp"
                                value={filterState.timestampStart}
                                onChange={(newValue) => {
                                    setFilterState({
                                        ...filterState,
                                        timestampStart: newValue,
                                        timestampEnd: null,
                                    });
                                }}
                                maxDate={filterState.timestampEnd}
                            />
                        </Grid>
                        <Grid item xs={12} md={4}>
                            <DateTimePicker
                                label="End Timestamp"
                                value={filterState.timestampEnd}
                                onChange={(newValue) => {
                                    setFilterState({
                                        ...filterState,
                                        timestampEnd: newValue,
                                    });
                                }}
                                minDate={filterState.timestampStart}
                            />
                        </Grid>
                    </LocalizationProvider>
                </Grid>
            </div>

            <div style={{
                marginTop: '100px',
            }}>
                <DataGrid
                    autoHeight
                    columns={columns}
                    rows={pageState.data}
                    pageSizeOptions={[10, 30, 50, 70, 100]}
                    rowCount={pageState.total}
                    loading={pageState.isLoading}
                    page={pageState.page}
                    pageSize={pageState.pageSize}
                    paginationMode="server"
                    paginationModel={pageState}
                    onPaginationModelChange={(newModel) => {
                        setPageState({
                            ...pageState,
                            page: newModel.page,
                            pageSize: newModel.pageSize,
                        });
                    }
                    }
                />
            </div>
        </div>
    )
}

export default LogTable;
