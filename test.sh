go run *.go reset < input.txt > output.txt
go run *.go --userID mickey createUser >> output.txt
go run *.go --userID mickey createUser >> output.txt
go run *.go --userID semsem createUser >> output.txt
go run *.go --userID hema createUser >> output.txt
go run *.go --userID mickey --group elshela --creator mickey --members semsem-hema --admins mickey  createGroup >> output.txt
go run *.go --userID yay createUser >> output.txt
go run *.go --userID nooo createUser >> output.txt
go run *.go --userID yay --group elshela --member yay --role member  joinGroup >> output.txt
go run *.go --userID nooo --group elshela --member nooo --role member  joinGroup >> output.txt
go run *.go --userID mickey --group elshela getRequests >> output.txt
