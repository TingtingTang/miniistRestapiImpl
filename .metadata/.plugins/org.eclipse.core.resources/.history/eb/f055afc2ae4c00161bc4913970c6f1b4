	package com.cn.ibm.miniist.tingtingtang.rest;

import java.sql.Timestamp;
import java.text.SimpleDateFormat;
import java.util.ArrayList;
import java.util.List;

import javax.ws.rs.Consumes;
import javax.ws.rs.DELETE;
import javax.ws.rs.FormParam;
import javax.ws.rs.GET;
import javax.ws.rs.POST;
import javax.ws.rs.PUT;
import javax.ws.rs.Path;
import javax.ws.rs.PathParam;
import javax.ws.rs.core.MediaType;
import javax.ws.rs.core.Response;
import javax.ws.rs.core.Response.Status;

import com.cn.ibm.miniist.tingtingtang.entities.MachineInfo;
import com.cn.ibm.miniist.tingtingtang.services.MachineInfoService;
import com.google.gson.Gson;

@Path("/")
public class RestAPI {
	@GET
	@Path("/hello")
	//@Produces("text/html")
	public Response helloMiniist() {
		
		String hello = "MiniIST";
		return Response.ok(hello).build();
	}

	/*
	 * This put is used to create a new machine info to the database
	 */
	
	@POST
	@Path("/test/{mName}")
	@Consumes(MediaType.APPLICATION_FORM_URLENCODED)
	public Response testDBConnection(
			@PathParam("mName") String mName)
	{
		MachineInfo machineInfo = new MachineInfo();
		MachineInfoService machineInfoService = new MachineInfoService();
		try 
		{
			if (mName != null) 
			{
				machineInfo = machineInfoService.findMachine(mName, "mlist");

				Gson gson = new Gson();
				String json = gson.toJson(machineInfo);
				return Response.status(201).entity(json).build();

			} 
			else
			{
				return Response.status(Status.BAD_REQUEST).build();
			}
		}
		catch (Exception e) 
		{
			return Response.status(Status.BAD_REQUEST).build();
		}
		
	}
	
	@PUT
	@Path("/create")
	@Consumes(MediaType.APPLICATION_FORM_URLENCODED)
	public Response createMachineInfo(
			@FormParam("userName") String userName,
			@FormParam("mName") String mName,
			@FormParam("hostName") String hostName,
			@FormParam("ipAddress") String ipAddress,
			@FormParam("description") String description)
	{
		MachineInfo machineInfo = new MachineInfo();
		MachineInfoService machineInfoService = new MachineInfoService();
		String createTime;
		SimpleDateFormat dataFormat = new SimpleDateFormat("yyyy-MM-dd HH:MM:SS");
		try{
			createTime = dataFormat.format(new Timestamp(new java.util.Date().getTime()));
			machineInfo.setUserName(userName);
			machineInfo.setMachineName(mName);
			machineInfo.setHostName(hostName);
			machineInfo.setIpAdd(ipAddress);
			machineInfo.setDes(description);
			machineInfo.setCreateTime(createTime);
			machineInfo.setUpdateTime(createTime);
			
			if(machineInfoService.createNewMachine(machineInfo, "mlist"))
			{
				return Response.status(201).entity("Succeed!").build();
			}
			else
			{
				return Response.status(Status.BAD_REQUEST).build();
			}
		}
		catch (Exception e) 
		{
			return Response.status(Status.BAD_REQUEST).build();
		}
	}
	
	/*
	 * This post is used to find the machine info according to the unique machine name
	 */
	@POST
	@Path("/find")
	@Consumes(MediaType.APPLICATION_FORM_URLENCODED)
	public Response findMachineInfo(
			@FormParam("mName") String mName
			) 
	{
		MachineInfo machineInfo = new MachineInfo();
		MachineInfoService machineInfoService = new MachineInfoService();
		try {
			if (mName != null) 
			{
				machineInfo = machineInfoService.findMachine(mName, "mlist");

				Gson gson = new Gson();
				String json = gson.toJson(machineInfo);
				
				return Response.status(201).entity(json).build();
			} 
			else
			{
				return Response.status(Status.BAD_REQUEST).build();
			}

		} 
		catch (Exception e) 
		{
			return Response.status(Status.BAD_REQUEST).build();
		}
	}
	
	/*
	 * This post is used to find the machine list under specific user
	 */
	@POST
	@Path("/find/mlist")
	@Consumes(MediaType.APPLICATION_FORM_URLENCODED)
	public Response findMachineList(
			@FormParam("userName") String userName
			) 
	{
		List<MachineInfo> machineList = new ArrayList<MachineInfo>();
		MachineInfoService machineInfoService = new MachineInfoService();
		try {
			if (userName != null) 
			{
				machineList = machineInfoService.findMachineList(userName, "mlist");

				Gson gson = new Gson();
				String json = gson.toJson(machineList.toArray());
				
				return Response.status(201).entity(json).build();
			} 
			else
			{
				return Response.status(Status.BAD_REQUEST).build();
			}

		} 
		catch (Exception e) 
		{
			return Response.status(Status.BAD_REQUEST).build();
		}
	}
	
	
	/*
	 * Update machine info from the database
	 */
	@POST
	@Path("/update")
	@Consumes(MediaType.APPLICATION_FORM_URLENCODED)
	public Response findMachineInfo(
			@FormParam("userName") String userName,
			@FormParam("mName") String mName,
			@FormParam("hostName") String hostName,
			@FormParam("ipAddress") String ipAddress,
			@FormParam("description") String description
			) 
	{
		MachineInfo machineInfo = new MachineInfo();
		MachineInfoService machineInfoService = new MachineInfoService();
		String updateTime;
		SimpleDateFormat dataFormat = new SimpleDateFormat("yyyy-MM-dd HH:MM:SS");
		try{
			updateTime = dataFormat.format(new Timestamp(new java.util.Date().getTime()));
			machineInfo.setUserName(userName);;
			machineInfo.setMachineName(mName);
			machineInfo.setHostName(hostName);
			machineInfo.setIpAdd(ipAddress);
			machineInfo.setDes(description);
			machineInfo.setUpdateTime(updateTime);
			
			if(machineInfoService.updateMachine(machineInfo, "mlist"))
			{
				return Response.status(201).entity("Succeed!").build();
			}
			else
			{
				return Response.status(Status.BAD_REQUEST).build();
			}
		}
		catch (Exception e) 
		{
			return Response.status(Status.BAD_REQUEST).build();
		}
	}
	/*
	 * Delete the machine from the database
	 */
	@DELETE
	@Path("/delete/{mName}")
	@Consumes(MediaType.APPLICATION_FORM_URLENCODED)
	public Response deleteMachineInfo(
			@PathParam("mName") String mName)
	{
		MachineInfoService machineInfoService = new MachineInfoService();
		try
		{
			if(machineInfoService.deleteMachine(mName, "mlist"))
			{
				return Response.status(200).entity("delete" + mName + "successfully !").build();
			}
			else
			{
				return Response.status(Status.BAD_REQUEST).entity(mName + "is not exist in the database!").build();
			}	
		}
		catch (Exception e)
		{
			return Response.status(Status.BAD_REQUEST).build();
		}
	}
}
