package repository

//go:generate sh -c "rm -rf mocks && mkdir -p mocks"
//go:generate minimock -i NoteRepository -o ./mocks/ -s "_minimock.go"
