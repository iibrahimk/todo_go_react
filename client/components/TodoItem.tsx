
import { Badge, Box, Flex, Spinner, Text } from "@chakra-ui/react";
import { FaCheckCircle } from "react-icons/fa";
import { MdDelete } from "react-icons/md";
import { Todo } from "./TodoList"
import { useMutation, useQueryClient } from "@tanstack/react-query";

const TodoItem = ({ todo }: { todo: Todo }) => {
	const queryClient = useQueryClient()

	const {mutate:updateTodo, isPending: isUpdating} = useMutation({
		mutationKey: ["updateTodo"],
		mutationFn: async ()=>{
			if(todo.completed) return alert("Task is already completed!")
			try{
				const res = await fetch(import.meta.env.SERVER_PORT+`http://127.0.0.1:4000/api/updatetodo/${todo.taskid}`,{
					method:"PATCH"
				})
				const data = await res.json()
				if(!res.ok){
					throw new Error(data.error || "Something went wrong!")
				}
				return data
			}
			catch(error){
				console.log(error)
			}
		},
		onSuccess: ()=>{
			queryClient.invalidateQueries({queryKey:["todos"]})
		}

	})
	const {mutate:deleteTodo, isPending:isDeleting} = useMutation({
		mutationKey:["deleteTodo"],
		mutationFn: async ()=>{
			try{
				const res = await fetch(import.meta.env.SERVER_PORT+`/deletetod/${todo.taskid}`,{
					method:"DELETE"
				})
				const data = await res.json()
				if(!res.ok){
					throw new Error(data.error || "Something went wrong!")
				}
				return data
			}
			catch(error){
				console.log(error)
			}
		},
		onSuccess:()=>{
			queryClient.invalidateQueries({queryKey:["todos"]})
		}
	})
	return (
		<Flex gap={2} alignItems={"center"}>
			<Flex
				flex={1}
				alignItems={"center"}
				border={"1px"}
				borderColor={"gray.600"}
				p={2}
				borderRadius={"lg"}
				justifyContent={"space-between"}
			>
				<Text
					color={todo.completed ? "green.200" : "orange.800"}
					textDecoration={todo.completed ? "line-through" : "none"}
				>
					{todo.body}
				</Text>
				{todo.completed && (
					<Badge ml='1' colorScheme='green'>
						Done
					</Badge>
				)}
				{!todo.completed && (
					<Badge ml='1' colorScheme='orange'>
						In Progress
					</Badge>
				)}
			</Flex>
			<Flex gap={2} alignItems={"center"}>
				<Box color={"green.500"} cursor={"pointer"} onClick={()=> updateTodo()}>
				{!isUpdating && <FaCheckCircle size={20}/>}
				{isUpdating && <Spinner size={"sm"}/>}
				</Box>
				<Box color={"red.500"} cursor={"pointer"} onClick={()=>deleteTodo()}>
					
					{!isDeleting && <MdDelete size={25} />}
					{isDeleting && <Spinner size={"sm"}/>}

				</Box>
			</Flex>
		</Flex>
	);
};
export default TodoItem;