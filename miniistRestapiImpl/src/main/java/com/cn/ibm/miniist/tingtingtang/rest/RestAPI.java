package com.cn.ibm.miniist.tingtingtang.rest;

import javax.ws.rs.FormParam;
import javax.ws.rs.GET;
import javax.ws.rs.POST;
import javax.ws.rs.Path;
import javax.ws.rs.Produces;
import javax.ws.rs.core.Response;

import com.cn.ibm.miniist.tingtingtang.entities.MachineInfo;
import com.cn.ibm.miniist.tingtingtang.services.MachineInfoService;
import com.google.gson.Gson;


@Path("/api")
public class RestAPI 
{	
	@GET
	@Produces("text/html")
	
	public Response helloMiniist()
	{
		return Response.ok("MiniIST").build();
	}
	
	//This post is used to find the machine info according to the unique machine name
	@POST
	@Path("/findMachineInfo")
	public Response findMachineInfo(
			@FormParam("mName")
			String mName
			)
	{
		MachineInfo machineInfo = new MachineInfo();
		MachineInfoService machineInfoService = new MachineInfoService();
		machineInfo = machineInfoService.findMachine(mName);
		Gson gson = new Gson();
		String json = gson.toJson(machineInfo);
		return Response.ok(json).build();
	}
	
	@POST
	@Path("/updateMachineInfo")
	public Response updateMachineInfo(
			@FormParam("mName")
			String mName,
			@FormParam("hostName")
			String hostName,
			@FormParam("ipAdd")
			String ipAdd,
			@FormParam("des")
			String des
			)
	{
		MachineInfo machineInfo = new MachineInfo();
		MachineInfoService machineInfoService = new MachineInfoService();
		
		machineInfo.setMachineName(mName);
		machineInfo.setHostName(hostName);
		machineInfo.setIpAdd(ipAdd);
		machineInfo.setDes(des);
		
		machineInfoService.updateMachine(mName, machineInfo);
		return Response.ok().build();
	}
}

