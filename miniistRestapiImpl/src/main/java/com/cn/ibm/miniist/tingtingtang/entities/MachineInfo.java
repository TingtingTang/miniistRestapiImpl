package com.cn.ibm.miniist.tingtingtang.entities;

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
	private String createTime;
	private String updateTime;
	
	/*public MachineInfo(String mName, String hostName, String ipAdd, int id, String des, Timestamp createTime, Timestamp updateTime)
	{
		this.mName = mName;
		this.hostName = hostName;
		this.ipAdd = ipAdd;
		this.id = id;
		this.des = des;
		this.createTime = createTime;
		this.updateTime = updateTime;
	}*/
	
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
	
	/*
	 * get and set ip address
	 */
	public String getIpAdd()
	{
		return ipAdd;
	}
	
	public void setIpAdd(String ipAdd)
	{
		this.ipAdd = ipAdd;
	}
	
	/*
	 * get and set id; but this id should be generated from MongoDB, 
	 * so this function need to be confirmed and reset later on.
	 */
	public int getId()
	{
		return id;
	}
	
	public void setId(int id)
	{
		this.id = id;
	}
	
	/*
	 * get and set description of machine
	 */
	public String getDes()
	{
		return des;
	}
	
	public void setDes(String des)
	{
		this.des = des;
	}
	
	/*
	 * get and set the create time of the machine.
	 */
	public String getCreateTime()
	{
		return createTime;
	}
	
	public void setCreateTime(String createTime)
	{
		this.createTime = createTime;
	}
	
	/*
	 * get and set the update time of the machine
	 */
	public String getUpdateTime()
	{
		return updateTime;
	}
	
	public void setUpdateTime(String updateTime)
	{
		this.updateTime = updateTime;
	}
}

