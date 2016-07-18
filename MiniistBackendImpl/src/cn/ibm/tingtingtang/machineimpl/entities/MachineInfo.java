package cn.ibm.tingtingtang.machineimpl.entities;

import java.sql.Timestamp;

/**
 * @author TingTing Tang
 * 
 * This entity represents machine information which includes machine name, host name, IP address, ID, description,
 * create time and update time of machine itself
 */

public class MachineInfo 
{
	private String mName;
	private String hostName;
	private String ipAdd;
	private int id;
	private String des;
	private Timestamp createTime;
	private Timestamp updateTime;
	
	public MachineInfo(String mName, String hostName, String ipAdd, int id, String des, Timestamp createTime, Timestamp updateTime)
	{
		this.mName = mName;
		this.hostName = hostName;
		this.ipAdd = ipAdd;
		this.id = id;
		this.des = des;
		this.createTime = createTime;
		this.updateTime = updateTime;
	}
	
	/*
	 * get and set machine name
	 */
	public String getMachineName()
	{
		return mName;
	}
	
	public void setMachineName(String mName)
	{
		this.mName = mName;
	}
	
	/*
	 * get and set host name
	 */
	public String getHostName()
	{
		return hostName;
	}
	
	public void setHostName(String hostName)
	{
		this.hostName = hostName;
	}
	
	
	
}
