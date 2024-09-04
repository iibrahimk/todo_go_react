
import { Button, Flex, Input, Spinner } from "@chakra-ui/react";
import { useMutation,useQueryClient } from "@tanstack/react-query";
import { useState } from "react";
import { IoMdAdd } from "react-icons/io";

const TodoForm = () => {
	const [newTodo, setNewTodo] = useState("");
	const queryClient = useQueryClient()
	const {mutate:createTodo, isPending: isAdding} = useMutation({
		mutationKey:["createTodo"],
		mutationFn: async ()=>{
			try{
				const res = await fetch(import.meta.env.SERVER_PORT+"/addtodo", {
					method: "POST",
					headers: {
					  "Content-Type": "application/json",
					},
					body: JSON.stringify({body: newTodo}), // Properly stringify the body
				  });
				const data = await res.json()
				console.log(data)
				if(res.status == 200){
					alert("Something Went Wrong!")
					console.log(data.Error)
					throw new Error( data.Error || "Something went wrong!")
				}
				setNewTodo("")
				return data
			}
			catch(error){
				throw new Error( "Something went wrong!")
			}
		},
		onSuccess:()=>{
			queryClient.invalidateQueries({queryKey:["todos"]})
		},	
	})
	return (
		<form onSubmit={()=>createTodo()}>
			<Flex gap={2}>
				<Input
					type='text'
					value={newTodo}
					onChange={(e) => setNewTodo(e.target.value)}
					ref={(input) => input && input.focus()}
				/>
				<Button
					mx={2}
					type='submit'
					_active={{
						transform: "scale(.97)",
					}}
				>
					{isAdding ? <Spinner size={"xs"} /> : <IoMdAdd size={30} />}
				</Button>
			</Flex>
		</form>
	);
};
export default TodoForm;