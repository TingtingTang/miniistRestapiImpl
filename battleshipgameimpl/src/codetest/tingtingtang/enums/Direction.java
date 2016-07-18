/**
 * 
 */
package codetest.tingtingtang.enums;

/**
 * @author TingTing Tang
 * 
 * The direction means the direction of each ship as being placed and it includes horizontal and vertical.
 */
public enum Direction 
{
	HORIZONTAL(0),
	VERTICAL(1);
	
	private int value;
	private Direction(int value)
	{
		this.value = value;
	}
	
	public int getValue()
	{
		return value;
	}
}
