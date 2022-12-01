NAME            =ethereumAccount
SRCS            = src/main.go src/keystore.go src/store.go
RM              = rm -f
FLAGS           = -Wall -Wextra -Werror -std=c++98

all:            $(NAME)

$(NAME):
	go mod tidy
	go build -o $(NAME) $(SRCS)

clean:
	$(RM) ./src/pass/* ./src/wallets_keys/*

fclean:         clean
	$(RM) $(NAME)

re:             fclean $(NAME)

.PHONY:         all clean fclean re