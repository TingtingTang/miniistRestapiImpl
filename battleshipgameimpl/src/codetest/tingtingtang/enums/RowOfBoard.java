/**
 * 
 */
package codetest.tingtingtang.enums;

/**
 * @author TingTing Tang
 * 
 * The enum for row of the board corresponds to the user's input which need to be changed to the row of the board,
 * so that the subsequent algorithms can be executed.
 */
public enum RowOfBoard 
{
	A(0), 
	B(1),
	C(2),
	D(3),
	E(4),
	F(5),
	G(6),
	H(7),
	I(8),
	J(9);
	
	private int value;
	private RowOfBoard(int value)
	{
		this.value = value;
	}
	
	public int getValue()
	{
		return value;
	}
}
