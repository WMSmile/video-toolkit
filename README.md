# Video Toolkit

A powerful video processing application built with Wails3 and Vue.js, providing a Windows-style desktop interface for audio and video conversion tasks.

## Features

- **Audio Conversion**: Convert video files to MP3 format
- **Batch Processing**: Process multiple files at once
- **ID3 Tag Editing**: Edit metadata for audio files
- **Theming Support**: Light, dark, and blue themes
- **Windows-style UI**: Familiar desktop interface

## Dependencies

### Required
- **Go 1.25+**: For backend development
- **Node.js 18+**: For frontend development
- **FFmpeg**: For audio/video conversion

### Optional
- **Wails3 CLI**: For development and building

## Getting Started

### 1. Install Dependencies

#### FFmpeg Installation
- **Windows**: Download from [FFmpeg官网](https://ffmpeg.org/download.html) and add to PATH
- **macOS**: `brew install ffmpeg`
- **Linux**: `sudo apt install ffmpeg` (Ubuntu/Debian)

#### Node.js Dependencies
```bash
cd frontend
npm install
```

### 2. Development Mode

Run the application in development mode with hot-reloading:

```bash
wails3 dev
```

### 3. Production Build

Build the application for production:

```bash
wails3 build
```

This will create a production-ready executable in the `bin` directory.

### 4. Generate Bindings

If you make changes to the Go backend, regenerate the frontend bindings:

```bash
wails3 generate bindings
```

## Project Structure

```
video-toolkit/
├── app/                 # Go backend services
│   ├── batch.go         # Batch processing service
│   ├── ffmpeg.go        # FFmpeg conversion service
│   └── id3.go           # ID3 tag editing service
├── frontend/            # Vue.js frontend
│   ├── src/             # Frontend source code
│   │   ├── views/       # Vue components
│   │   ├── router/      # Vue Router configuration
│   │   └── styles/      # CSS/SCSS styles
│   ├── bindings/        # Auto-generated Wails bindings
│   └── package.json     # Frontend dependencies
├── main.go              # Application entry point
└── README.md            # This documentation
```

## Usage

1. **Audio Conversion**: Select a video file and convert it to MP3 format
2. **Batch Processing**: Select multiple files for batch conversion
3. **ID3 Tag Editing**: Edit metadata like title, artist, and album for audio files
4. **Theme Switching**: Change between light, dark, and blue themes in settings

## Backend Services

### FFmpeg Service
- **ToMP3(input string)**: Convert video file to MP3 format

### ID3 Service
- **Read(file string)**: Read ID3 tags from audio file
- **Write(file string, tags Tags)**: Write ID3 tags to audio file

### Batch Service
- **BatchConvert(inputs []string)**: Convert multiple files to MP3

## Frontend Components

- **Home**: Main dashboard
- **About**: Application information
- **Settings**: Theme and configuration settings

## Troubleshooting

### Common Issues

1. **FFmpeg not found**: Ensure FFmpeg is installed and added to PATH
2. **Bindings not updated**: Run `wails3 generate bindings` after backend changes
3. **Frontend changes not visible**: Check if development server is running with `wails3 dev`

### Development Tips

- Use `wails3 dev` for real-time updates during development
- Check the browser console for frontend errors
- Check the terminal for backend errors

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License.

## Acknowledgements

- [Wails3](https://wails.io/) - Go + Web frontend framework
- [Vue.js](https://vuejs.org/) - JavaScript framework
- [Element Plus](https://element-plus.org/) - UI library
- [FFmpeg](https://ffmpeg.org/) - Audio/video processing
- [bogem/id3v2](https://github.com/bogem/id3v2) - ID3 tag library

Happy coding!