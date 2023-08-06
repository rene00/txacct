"""empty message

Revision ID: 36a7759d44e8
Revises: 
Create Date: 2023-07-30 16:29:29.470275

"""
from alembic import op
import sqlalchemy as sa


# revision identifiers, used by Alembic.
revision = '36a7759d44e8'
down_revision = None
branch_labels = None
depends_on = None


def upgrade():
    # ### commands auto generated by Alembic - please adjust! ###
    op.create_table('sa3',
    sa.Column('id', sa.Integer(), nullable=False),
    sa.Column('code', sa.Integer(), nullable=False),
    sa.Column('name', sa.String(), nullable=False),
    sa.PrimaryKeyConstraint('id'),
    sa.UniqueConstraint('code', 'name')
    )
    op.create_table('sa4',
    sa.Column('id', sa.Integer(), nullable=False),
    sa.Column('code', sa.Integer(), nullable=False),
    sa.Column('name', sa.String(), nullable=False),
    sa.PrimaryKeyConstraint('id'),
    sa.UniqueConstraint('code', 'name')
    )
    op.create_table('state',
    sa.Column('id', sa.Integer(), nullable=False),
    sa.Column('name', sa.String(), nullable=False),
    sa.PrimaryKeyConstraint('id'),
    sa.UniqueConstraint('name')
    )
    op.create_table('transaction',
    sa.Column('id', sa.Integer(), nullable=False),
    sa.Column('memo', sa.String(), nullable=False),
    sa.PrimaryKeyConstraint('id')
    )
    op.create_table('postcode',
    sa.Column('id', sa.Integer(), nullable=False),
    sa.Column('postcode', sa.String(), nullable=False),
    sa.Column('locality', sa.String(), nullable=False),
    sa.Column('state_id', sa.Integer(), nullable=False),
    sa.Column('sa3_id', sa.Integer(), nullable=True),
    sa.Column('sa4_id', sa.Integer(), nullable=True),
    sa.ForeignKeyConstraint(['sa3_id'], ['sa3.id'], ),
    sa.ForeignKeyConstraint(['sa4_id'], ['sa4.id'], ),
    sa.ForeignKeyConstraint(['state_id'], ['state.id'], ),
    sa.PrimaryKeyConstraint('id'),
    sa.UniqueConstraint('postcode', 'locality', name='postcode_locality')
    )
    # ### end Alembic commands ###


def downgrade():
    # ### commands auto generated by Alembic - please adjust! ###
    op.drop_table('postcode')
    op.drop_table('transaction')
    op.drop_table('state')
    op.drop_table('sa4')
    op.drop_table('sa3')
    # ### end Alembic commands ###
