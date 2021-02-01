mockery --dir ./filewalker --name dirReader --output ./filewalker/mocks --filename dirReader_mock.go --structname TestDirReader 
mockery --dir ./filewalker --name filesProvider --output ./filewalker/mocks --filename filesProvider_mock.go --structname TestFilesProvider 
mockery --dir ./filewalker --name fileHandler --output ./filewalker/mocks --filename filesHandler_mock.go --structname TestFileHandler 

