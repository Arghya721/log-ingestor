import React, { useState, useEffect } from 'react';
import { Navbar } from '../components/Navbar';
import { DataGrid } from '@mui/x-data-grid';
import { Container } from '@mui/material';


const LogTable = () => {

    const [pageState, setPageState] = useState({
        isLoading: true,
        data: [],
        total: 0,
        page: 0,
        pageSize: 10,
    });

    const columns = [
        { field: 'id', headerName: 'ID', width: 100 },
        { field: 'level', headerName: 'Level', width: 130 },
        { field: 'message', headerName: 'Message', width: 400 },
        { field: 'resourceId' , headerName: 'Resource ID', width: 200 },
        { field: 'timestamp', headerName: 'Timestamp', width: 200 },
        { field: 'traceId', headerName: 'Trace ID', width: 200 },
        { field: 'spanId', headerName: 'Span ID', width: 200 },
        { field: 'commit' , headerName: 'Commit', width: 200 },
        { field: 'parentResourceId', headerName: 'Parent Resource ID', width: 200 },
    ];

    useEffect(() => {
        const fetchLogs = async () => {
            setPageState({
                ...pageState,
                isLoading: true,
            });
            const res = await fetch(`http://localhost:1323/internal/logs?page=${pageState.page}&limit=${pageState.pageSize}`);
            const data = await res.json();
            setPageState({
                ...pageState,
                isLoading: false,
                data: data.logs,
                total: data.total,
            });
        }
        fetchLogs();

    }, [pageState.page, pageState.pageSize]);




    return (
        <div>
            <Navbar
                title="Logs Visualization"
            />
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
