import React, { useState, useEffect } from 'react';
import { Navbar } from '../components/Navbar';
import { DataGrid } from '@mui/x-data-grid';
import { Container } from '@mui/material';


const LogTable = () => {

    const [pageState, setPageState] = useState({
        isLoading: true,
        data: [],
        total: 0,
        page: 1,
        pageSize: 10,

    });






    return (
        <div>
            <Navbar
                title="Logs Visualization"
            />
            <Container style={{ marginTop: 100 , marginBottom: 100 }}>
                <DataGrid
                    rows={pageState.data}
                    {...data}
                    rowCount={pageState.total}
                    loading={pageState.isLoading}
                    pageSizeOptions={[5]}
                    paginationModel={paginationModel}
                    paginationMode="server"
                    onPaginationModelChange={setPaginationModel}
                />
            </Container>
        </div>
    )
}

export default LogTable;
